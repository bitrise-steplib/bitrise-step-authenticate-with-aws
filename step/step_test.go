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
	tests := []struct {
		name           string
		env            map[string]string
		expectedConfig Config
		expectError    bool
	}{
		{
			name: "valid config with account key",
			env: map[string]string{
				"build_url":         "build-url",
				"build_api_token":   "build-token",
				"access_key_id":     "access-key-id",
				"secret_access_key": "secret-access-key",
				"audience":          "",
				"role_arn":          "",
				"session_name":      "",
				"region":            "us-east-1",
				"docker_login":      "false",
				"verbose":           "false",
			},
			expectedConfig: Config{
				BuildURL:        "build-url",
				BuildToken:      "build-token",
				AccessKeyId:     "access-key-id",
				SecretAccessKey: "secret-access-key",
				Audience:        "",
				RoleArn:         "",
				Region:          "us-east-1",
				SessionName:     "",
				DockerLogin:     false,
			},
			expectError: false,
		},
		{
			name: "valid config with identity token",
			env: map[string]string{
				"build_url":         "build-url",
				"build_api_token":   "build-token",
				"access_key_id":     "",
				"secret_access_key": "",
				"audience":          "audience",
				"role_arn":          "role-arn",
				"session_name":      "session-name",
				"region":            "us-east-1",
				"docker_login":      "false",
				"verbose":           "false",
			},
			expectedConfig: Config{
				BuildURL:        "build-url",
				BuildToken:      "build-token",
				AccessKeyId:     "",
				SecretAccessKey: "",
				Audience:        "audience",
				RoleArn:         "role-arn",
				Region:          "us-east-1",
				SessionName:     "session-name",
				DockerLogin:     false,
			},
			expectError: false,
		},
		{
			name: "error when both account key and identity config are set",
			env: map[string]string{
				"build_url":         "build-url",
				"build_api_token":   "build-token",
				"access_key_id":     "access-key-id",
				"secret_access_key": "secret-access-key",
				"audience":          "audience",
				"role_arn":          "role-arn",
				"session_name":      "session-name",
				"region":            "us-east-1",
				"docker_login":      "false",
				"verbose":           "false",
			},
			expectedConfig: Config{},
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEnvRepository := mocks.NewRepository(t)
			for k, v := range tt.env {
				mockEnvRepository.On("Get", k).Return(v)
			}

			inputParser := stepconf.NewInputParser(mockEnvRepository)
			mockFactory := mocks.NewFactory(t)
			exporter := export.NewExporter(mocks.NewFactory(t))
			sut := NewAuthenticator(inputParser, mockEnvRepository, mockFactory, exporter, log.NewLogger())

			receivedConfig, err := sut.ProcessConfig()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedConfig, receivedConfig)
			}

			mockEnvRepository.AssertExpectations(t)
		})
	}
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
