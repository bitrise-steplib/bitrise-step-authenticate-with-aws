package step

import (
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-steplib/bitrise-step-get-identity-token/api"
)

func (a Authenticator) identityToken(url string, token stepconf.Secret, aud string) (string, error) {
	client := api.NewDefaultAPIClient(url, token, a.logger)

	parameter := api.GetIdentityTokenParameter{
		Audience: aud,
	}
	response, err := client.GetIdentityToken(parameter)
	if err != nil {
		return "", err
	}

	return response.Token, nil
}
