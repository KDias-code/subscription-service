package configs

import (
	"time"
)

type Config struct {
	Http     HttpConfig
	Log      LogConfig
	Postgres PostgresConfig
}

type HttpConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type PostgresConfig struct {
	Dsn string
}

type LogConfig struct {
	Level int32
	Json  bool
}
