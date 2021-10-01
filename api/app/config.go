package app

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ApplicationKey        string        `envconfig:"APPLICATION_KEY"`
	AuthMocked            bool          `envconfig:"AUTH_MOCKED" default:"false"`
	DBUser                string        `envconfig:"DB_USER"`
	DBPass                string        `envconfig:"DB_PASS"`
	DBName                string        `envconfig:"DB_NAME"`
	DBHost                string        `envconfig:"DB_HOST"`
	DBSSLMode             string        `envconfig:"DB_SSL_MODE"`
	DBPoolMaxConns        int           `envconfig:"DB_POOL_MAX_CONNS" default:"10"`
	DBPoolMaxConnIdleTime time.Duration `envconfig:"DB_POOL_MAX_CONN_IDLE_TIME" default:"30m"`
	DBPoolMinConns        int           `envconfig:"DB_POOL_MIN_CONNS" default:"5"`
	AWSS3Endpoint         string        `envconfig:"AWS_S3_ENDPOINT"`
	AWSS3Region           string        `envconfig:"AWS_REGION"`
	AWSS3DisableSSL       bool          `envconfig:"AWS_S3_DISABLE_SSL"`
	AWSS3ForcePathStyle   bool          `envconfig:"AWS_S3_FORCE_PATH_STYLE"`
	AWSS3Bucket           string        `envconfig:"AWS_S3_BUCKET"`
}

// GetConfig returns environment variable config
func GetConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("water", &cfg); err != nil {
		log.Fatal(err.Error())
	}
	return &cfg, nil
}

// AWSConfig
func AWSConfig(cfg Config) aws.Config {
	s3Config := aws.NewConfig().WithCredentials(credentials.NewEnvCredentials())
	s3Config.WithDisableSSL(cfg.AWSS3DisableSSL)
	s3Config.WithS3ForcePathStyle(cfg.AWSS3ForcePathStyle)
	s3Config.WithRegion(cfg.AWSS3Region)
	if cfg.AWSS3Endpoint != "" {
		s3Config.WithEndpoint(cfg.AWSS3Endpoint)
	}
	return *s3Config
}
