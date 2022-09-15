package main

import "github.com/caarlos0/env/v6"

var config Config

type Config struct {
	Port      uint   `env:"PORT" envDefault:"8545"`
	NodeURL   string `env:"NODE_URL,required"`
	NodeWSURL string `env:"NODE_WS_URL"`
}

// LoadConfig populates the config
func LoadConfig() error {
	if err := env.Parse(&config); err != nil {
		return err
	}
	if config.NodeWSURL == "" {
		config.NodeWSURL = handleInfuraNodeWS(config.NodeURL)
	}
	return nil
}
