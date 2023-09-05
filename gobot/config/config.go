package config

import (
	"log"
	"os"
)

type Config struct {
	Token string
}

func NewDefaultConfig() *Config {
	token, found := os.LookupEnv("TOKEN")
	if !found {
		log.Fatal("TOKEN environment variable is not set")
	}
	return &Config{
		Token: token,
	}
}

func (c *Config) GetToken() string {
	return c.Token
}
