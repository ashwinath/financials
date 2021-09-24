package config

import "time"

type Config struct {
	Server   Server
	Database Database
}

type Server struct {
	Port                  int           `validate:"required"`
	WriteTimeoutInSeconds time.Duration `validate:"required"`
	ReadTimeoutInSeconds  time.Duration `validate:"required"`
}

type Database struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
	Name     string `validate:"required"`
	TimeZone string `validate:"required"`
}
