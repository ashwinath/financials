package config

import "time"

type Config struct {
	Server Server
}

type Server struct {
	Port                  int           `validate:"required"`
	WriteTimeoutInSeconds time.Duration `validate:"required"`
	ReadTimeoutInSeconds  time.Duration `validate:"required"`
}
