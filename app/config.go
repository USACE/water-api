package app

import (
	"time"
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
}
