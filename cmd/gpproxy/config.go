package main

import "github.com/caarlos0/env/v6"

var config Config

type Config struct {
	Port    uint   `env:"PORT" envDefault:"8545"`
	NodeURL string `env:"NODE_URL,required"`
}

// LoadConfig populates the config
func LoadConfig() error {
	err := env.Parse(&config)
	return err
}
