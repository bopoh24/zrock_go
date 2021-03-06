package settings

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// App global apps settings
var App struct {
	BindAdd      string `envconfig:"BIND_ADDR" default:":8080"`
	LogLevel     string `envconfig:"LOG_LEVEL" default:"debug"`
	DatabaseURL  string `envconfig:"DATABASE_URL" default:"host=localhost user=postgres dbname=zrock_api_dev sslmode=disable"`
	TokenSecret  string `envconfig:"TOKEN_SECRET" default:"my_token_secret_string"`
	SMTPUser     string `envconfig:"SMTP_USER"`
	SMTPPassword string `envconfig:"SMTP_PASSWORD"`
	SMTPHost     string `envconfig:"SMTP_HOST"`
	SMTPPort     int    `envconfig:"SMTP_PORT"`
	SMTPFrom     string `envconfig:"SMTP_FROM"`
}

func init() {
	err := envconfig.Process("ZROCK", &App)
	if err != nil {
		log.Fatal("Config initialization failed. ERROR:", err.Error())
	}
}
