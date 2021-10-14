package config

import "time"

// Config is the base config entrypoint
type Config struct {
	Server                       Server        `validate:"required,dive"`
	Database                     Database      `validate:"required,dive"`
	AlphaVantageAPIKey           string        `validate:"required"`
	PortfolioCalculationInterval time.Duration `validate:"required"`
}

// Server contains the server related configurations
type Server struct {
	Port                  int           `validate:"required"`
	WriteTimeoutInSeconds time.Duration `validate:"required"`
	ReadTimeoutInSeconds  time.Duration `validate:"required"`
	ReactFilePath string `validate:"required"`
}

// Database contains the database related configurations
type Database struct {
	Host            string `validate:"required"`
	Port            int    `validate:"required"`
	User            string `validate:"required"`
	Password        string `validate:"required"`
	Name            string `validate:"required"`
	TimeZone        string `validate:"required"`
	BatchInsertSize int    `validate:"required"`
}
