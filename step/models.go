package step

import (
	"github.com/bitrise-io/go-steputils/v2/stepconf"
)

type Input struct {
	BuildURL        string          `env:"build_url,required"`
	BuildToken      stepconf.Secret `env:"build_api_token,required"`
	AccessKeyId     stepconf.Secret `env:"access_key_id"`
	SecretAccessKey stepconf.Secret `env:"secret_access_key"`
	Audience        string          `env:"audience"`
	RoleArn         string          `env:"role_arn"`
	Region          string          `env:"region"`
	SessionName     string          `env:"session_name"`
	DockerLogin     bool            `env:"docker_login,opt[true,false]"`
	Verbose         bool            `env:"verbose,opt[true,false]"`
}

type Config struct {
	BuildURL        string
	BuildToken      stepconf.Secret
	AccessKeyId     stepconf.Secret
	SecretAccessKey stepconf.Secret
	Audience        string
	RoleArn         string
	Region          string
	SessionName     string
	DockerLogin     bool
}

type Result struct {
	AccessKeyId     string
	SecretAccessKey string
	SessionToken    string
}
