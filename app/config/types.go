package config

// Config is the base config entrypoint
type Config struct {
	Database           Database `validate:"required,dive"`
	AlphaVantageAPIKey string   `validate:"required"`
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
