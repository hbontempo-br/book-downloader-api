package config

import (
	"github.com/plaid/go-envvar/envvar"
	"go.uber.org/zap"
)

type DBConfig struct {
	Address  string `envvar:"ADDRESS"`
	Port     int    `envvar:"PORT"`
	DBName   string `envvar:"NAME"`
	User     string `envvar:"USER"`
	Password string `envvar:"PASSWORD"`
}

type MinioConfig struct {
	Endpoint  string `envvar:"ENDPOINT"`
	AccessKey string `envvar:"ACCESS_KEY"`
	SecretKey string `envvar:"SECRET_KEY"`
	SSL       bool   `envvar:"SSL"`
}

type BucketConfig struct {
	Name                    string `envvar:"NAME"`
	DefaultDownloadLinkTime int    `envvar:"DEFAULT_DOWNLOAD_LINK_TIME" default:"60"`
}

type EnvVars struct {
	DBConfig     DBConfig     `envvar:"DB_"`
	MinioConfig  MinioConfig  `envvar:"MINIO_"`
	BucketConfig BucketConfig `envvar:"BUCKET_"`
	ServerPort   int          `envvar:"PORT" default:"3000"`
	Environment  string       `envvar:"ENVIRONMENT" default:"local"`
}

var envConfig *EnvVars

func LoadEnvVars() (*EnvVars, error) {

	// TODO: Add custom errors to package
	if envConfig != nil {
		zap.S().Debug("Environments variables already loaded, reaching value on memory")
		return envConfig, nil
	}
	var newEnvConfig EnvVars
	err := envvar.Parse(&newEnvConfig)
	if err != nil {
		zap.S().Debug("Unable to load environments variables correctly", err)
		return nil, err
	}
	envConfig = &newEnvConfig
	zap.S().Debug("Environments variables successfully loaded")
	return &newEnvConfig, nil
}
