package step

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/go-utils/v2/command"
)

func (a Authenticator) authenticate(cfg aws.Config, roleArn, sessionName, identityToken string) (Result, error) {
	ctx := context.Background()
	client := sts.NewFromConfig(cfg)
	response, err := client.AssumeRoleWithWebIdentity(ctx, &sts.AssumeRoleWithWebIdentityInput{
		RoleArn:          aws.String(roleArn),
		RoleSessionName:  aws.String(sessionName),
		WebIdentityToken: aws.String(identityToken),
	})
	if err != nil {
		return Result{}, fmt.Errorf("failed to assume role with web identity: %w", err)
	}

	result := Result{}

	if response.Credentials.AccessKeyId == nil {
		return Result{}, errors.New("AWS access key id is empty")
	}
	result.AccessKeyId = *response.Credentials.AccessKeyId

	if response.Credentials.SecretAccessKey == nil {
		return Result{}, errors.New("AWS secret access key is empty")
	}
	result.SecretAccessKey = *response.Credentials.SecretAccessKey

	if response.Credentials.SessionToken == nil {
		return Result{}, errors.New("AWS session token is empty")
	}
	result.SessionToken = *response.Credentials.SessionToken

	return result, nil
}

func (a Authenticator) loginWithDocker(cfg aws.Config, result Result) error {
	ctx := context.Background()
	cfg.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		aws.ToString(&result.AccessKeyId),
		aws.ToString(&result.SecretAccessKey),
		aws.ToString(&result.SessionToken),
	))
	client := ecr.NewFromConfig(cfg)

	auth, err := client.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return err
	}
	if len(auth.AuthorizationData) == 0 || auth.AuthorizationData[0].AuthorizationToken == nil {
		return fmt.Errorf("no authorization data returned")
	}

	decoded, err := base64.StdEncoding.DecodeString(*auth.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		return err
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("unexpected token format: %q", string(decoded))
	}

	password := parts[1]
	registry := strings.TrimPrefix(aws.ToString(auth.AuthorizationData[0].ProxyEndpoint), "https://")

	cmd := a.commandFactory.Create("docker", []string{"login", "--username", "AWS", "--password-stdin", registry}, &command.Opts{
		Stdin: strings.NewReader(password),
	})
	if output, err := cmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
		if errorutil.IsExitStatusError(err) {
			a.logger.Errorf("Docker login output: %s", output)
		}
		return err
	}

	return nil
}
