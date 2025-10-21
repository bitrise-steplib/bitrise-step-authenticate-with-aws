package step

import (
	"github.com/bitrise-io/go-steputils/v2/stepconf"
)

type Input struct {
	BuildURL    string          `env:"build_url,required"`
	BuildToken  stepconf.Secret `env:"build_api_token,required"`
	Audience    string          `env:"audience,required"`
	RoleArn     string          `env:"role_arn,required"`
	Region      string          `env:"region,required"`
	SessionName string          `env:"session_name"`
	DockerLogin bool            `env:"docker_login,opt[true,false]"`
	Verbose     bool            `env:"verbose,opt[true,false]"`
}

type Config struct {
	BuildURL    string
	BuildToken  stepconf.Secret
	Audience    string
	RoleArn     string
	Region      string
	SessionName string
	DockerLogin bool
}

type Result struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}
