package config

import (
	"flag"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type configFile string

func (c *configFile) String() string {
	return string(*c)
}

func (c *configFile) Set(value string) error {
	*c = configFile(value)
	return nil
}

// Load loads the config file
func Load() (*Config, error) {
	var c configFile
	flag.Var(&c, "config", "Path to a configuration file.")
	flag.Parse()

	if c == "" {
		return nil, fmt.Errorf("must set config file with -config flag")
	}

	viper.SetConfigFile(c.String())
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	if err := validator.New().Struct(config); err != nil {
		return nil, err
	}

	return &config, nil
}
