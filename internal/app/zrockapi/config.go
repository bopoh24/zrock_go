package zrockapi

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config contains server configuration
type Config struct {
	BindAdd     string `envconfig:"BIND_ADDR" default:":8080"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"debug"`
	DatabaseURL string `envconfig:"DATABASE_URL" default:"host=localhost user=postgres dbname=zrock_api_dev sslmode=disable"`
}

// NewConfig returns Config
func NewConfig() *Config {
	config := &Config{}
	// Reading config from env
	err := envconfig.Process("ZROCK", config)
	if err != nil {
		log.Fatal("Config initialization failed. ERROR:", err.Error())
	}
	return config
}
