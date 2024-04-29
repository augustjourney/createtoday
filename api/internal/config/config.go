package config

import (
	"flag"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	BaseURL              string        `env:"BASE_URL"`
	ServerAddress        string        `env:"SERVER_ADDRESS"`
	DatabaseDSN          string        `env:"DATABASE_DSN"`
	JWT_TOKEN_EXP        time.Duration `env:"JWT_TOKEN_EXP"`
	JWT_TOKEN_SECRET_KEY string        `env:"JWT_TOKEN_SECRET_KEY"`
	JWT_SIGNING_METHOD   jwt.SigningMethod
}

var config Config

func New() *Config {
	var flagServerAddress = flag.String("a", "localhost:8080", "Server address on which server is running")
	var flagDatabaseDSN = flag.String("d", "", "Database DSN")

	config.JWT_SIGNING_METHOD = jwt.SigningMethodHS256
	config.JWT_TOKEN_SECRET_KEY = "super-secret-key"
	config.JWT_TOKEN_EXP = time.Second * 30
	config.ServerAddress = *flagServerAddress

	if databaseDSN := os.Getenv("DATABASE_DSN"); databaseDSN != "" {
		config.DatabaseDSN = databaseDSN
	}

	if *flagDatabaseDSN != "" {
		config.DatabaseDSN = *flagDatabaseDSN
	}

	return &config

}
