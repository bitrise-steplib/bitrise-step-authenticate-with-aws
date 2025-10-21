package step

import (
	"testing"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-authenticate-with-aws/step/mocks"
	"github.com/stretchr/testify/assert"
)

func TestConfigParsing(t *testing.T) {
	config := Config{
		BuildURL:    "build-url",
		BuildToken:  stepconf.Secret("build-token"),
		Audience:    "audience",
		RoleArn:     "role-arn",
		Region:      "region",
		SessionName: "session-name",
		DockerLogin: true,
	}

	mockEnvRepository := mocks.NewRepository(t)
	mockEnvRepository.On("Get", "build_url").Return(config.BuildURL)
	mockEnvRepository.On("Get", "build_api_token").Return(string(config.BuildToken))
	mockEnvRepository.On("Get", "audience").Return(config.Audience)
	mockEnvRepository.On("Get", "role_arn").Return(config.RoleArn)
	mockEnvRepository.On("Get", "region").Return(config.Region)
	mockEnvRepository.On("Get", "session_name").Return(config.SessionName)
	mockEnvRepository.On("Get", "docker_login").Return("true")
	mockEnvRepository.On("Get", "verbose").Return("false")

	inputParser := stepconf.NewInputParser(mockEnvRepository)
	mockFactory := mocks.NewFactory(t)
	exporter := export.NewExporter(mocks.NewFactory(t))
	sut := NewAuthenticator(inputParser, mockEnvRepository, mockFactory, exporter, log.NewLogger())

	receivedConfig, err := sut.ProcessConfig()
	assert.NoError(t, err)
	assert.Equal(t, config, receivedConfig)

	mockEnvRepository.AssertExpectations(t)
}

func TestExport(t *testing.T) {
	result := Result{
		AccessKeyId:     "access-key-id",
		SecretAccessKey: "secret-access-key",
		SessionToken:    "session-token",
	}

	mockFactory := mocks.NewFactory(t)
	mockFactory.On("Create", "envman", mockParameters("AWS_ACCESS_KEY_ID", result.AccessKeyId), (*command.Opts)(nil)).Return(testCommand())
	mockFactory.On("Create", "envman", mockParameters("AWS_SECRET_ACCESS_KEY", result.SecretAccessKey), (*command.Opts)(nil)).Return(testCommand())
	mockFactory.On("Create", "envman", mockParameters("AWS_SESSION_TOKEN", result.SessionToken), (*command.Opts)(nil)).Return(testCommand())

	mockEnvRepository := mocks.NewRepository(t)
	inputParser := stepconf.NewInputParser(mockEnvRepository)
	exporter := export.NewExporter(mockFactory)
	sut := NewAuthenticator(inputParser, mockEnvRepository, mockFactory, exporter, log.NewLogger())

	err := sut.Export(result)
	assert.NoError(t, err)

	mockEnvRepository.AssertExpectations(t)
}

func testCommand() command.Command {
	factory := command.NewFactory(env.NewRepository())
	return factory.Create("pwd", []string{}, nil)
}

func mockParameters(key, value string) []string {
	return []string{"add", "--key", key, "--value", value, "--sensitive"}
}
