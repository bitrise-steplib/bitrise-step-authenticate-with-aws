package step

import (
	"context"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

const (
	idKey     = "AWS_ACCESS_KEY_ID"
	secretKey = "AWS_SECRET_ACCESS_KEY"
	tokenKey  = "AWS_SESSION_TOKEN"
)

type Authenticator struct {
	inputParser    stepconf.InputParser
	envRepository  env.Repository
	commandFactory command.Factory
	exporter       export.Exporter
	logger         log.Logger
}

func NewAuthenticator(inputParser stepconf.InputParser, envRepository env.Repository, commandFactory command.Factory, exporter export.Exporter, logger log.Logger) Authenticator {
	return Authenticator{
		inputParser:    inputParser,
		envRepository:  envRepository,
		commandFactory: commandFactory,
		exporter:       exporter,
		logger:         logger,
	}
}

func (a Authenticator) ProcessConfig() (Config, error) {
	var input Input
	err := a.inputParser.Parse(&input)
	if err != nil {
		return Config{}, err
	}

	stepconf.Print(input)
	a.logger.Println()
	a.logger.EnableDebugLog(input.Verbose)

	accessKeyAuth := input.AccessKeyId != "" && input.SecretAccessKey != ""
	identityAuth := input.Audience != "" && input.RoleArn != "" && input.Region != ""

	if accessKeyAuth && identityAuth {
		return Config{}, fmt.Errorf("only one authentication method can be used at a time (either Access Key or Identity Token)")
	}

	if !accessKeyAuth && !identityAuth {
		return Config{}, fmt.Errorf("no valid authentication method set (provide Access Key or Identity Token details)")
	}

	return Config{
		BuildURL:        input.BuildURL,
		BuildToken:      input.BuildToken,
		AccessKeyId:     input.AccessKeyId,
		SecretAccessKey: input.SecretAccessKey,
		Audience:        input.Audience,
		RoleArn:         input.RoleArn,
		Region:          input.Region,
		SessionName:     input.SessionName,
		DockerLogin:     input.DockerLogin,
	}, nil
}

func (a Authenticator) Run(config Config) (Result, error) {
	ctx := context.Background()
	awscfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(config.Region))
	if err != nil {
		return Result{}, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	var result Result

	if config.AccessKeyId != "" && config.SecretAccessKey != "" {
		result = Result{
			AccessKeyId:     string(config.AccessKeyId),
			SecretAccessKey: string(config.SecretAccessKey),
			SessionToken:    "",
		}
	} else {
		identityToken, err := a.identityToken(config.BuildURL, config.BuildToken, config.Audience)
		if err != nil {
			return Result{}, err
		}

		a.logger.Printf("Identity token fetched.\n")

		result, err = a.authenticate(awscfg, config.RoleArn, config.SessionName, identityToken)
		if err != nil {
			return Result{}, fmt.Errorf("AWS authenticate failure: %w", err)
		}
	}

	a.logger.Printf("Successful AWS authentication.\n")

	if config.DockerLogin {
		if err := a.loginWithDocker(awscfg, result); err != nil {
			return Result{}, fmt.Errorf("docker login failure: %w", err)
		}

		a.logger.Printf("Successful Docker login.\n")
	}

	return result, nil
}

func (a Authenticator) Export(result Result) error {
	a.logger.Printf("The following outputs are exported as environment variables:")

	values := map[string]string{
		idKey:     result.AccessKeyId,
		secretKey: result.SecretAccessKey,
		tokenKey:  result.SessionToken,
	}

	for key, value := range values {
		err := a.exporter.ExportSecretOutput(key, value)
		if err != nil {
			return err
		}

		a.logger.Donef("$%s", key)
	}

	return nil
}
