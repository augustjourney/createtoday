package config

import (
	"createtodayapi/internal/logger"
	"flag"
	"github.com/joho/godotenv"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	BaseURL            string        `env:"BASE_URL"`
	ServerAddress      string        `env:"SERVER_ADDRESS"`
	DatabaseDSN        string        `env:"DATABASE_DSN"`
	JwtTokenExp        time.Duration `env:"JWT_TOKEN_EXP"`
	MagicLinkExp       time.Duration `env:"MAGIC_LINK_EXP"`
	JwtTokenSecretKey  string        `env:"JWT_TOKEN_SECRET_KEY"`
	JwtSigningMethod   jwt.SigningMethod
	HeroAppBaseURL     string `env:"HERO_APP_BASE_URL"`
	AwsSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY"`
	AwsAccessKeyId     string `env:"AWS_ACCESS_KEY_ID"`
	AwsRegion          string `env:"AWS_REGION"`
}

var config Config

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Error("No .env file found")
	}
}

func New() *Config {
	var flagServerAddress = flag.String("a", "localhost:8080", "Server address on which server is running")
	var flagDatabaseDSN = flag.String("d", "", "Database DSN")

	config.JwtSigningMethod = jwt.SigningMethodHS256
	config.JwtTokenExp = time.Hour * 720
	config.MagicLinkExp = time.Minute * 1
	config.ServerAddress = *flagServerAddress

	if databaseDSN := os.Getenv("DATABASE_DSN"); databaseDSN != "" {
		config.DatabaseDSN = databaseDSN
	}

	if *flagDatabaseDSN != "" {
		config.DatabaseDSN = *flagDatabaseDSN
	}

	if JwtTokenSecretKey := os.Getenv("JWT_TOKEN_SECRET_KEY"); JwtTokenSecretKey != "" {
		config.JwtTokenSecretKey = JwtTokenSecretKey
	}

	if HeroAppBaseURL := os.Getenv("HERO_APP_BASE_URL"); HeroAppBaseURL != "" {
		config.HeroAppBaseURL = HeroAppBaseURL
	}

	if AwsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY"); AwsSecretAccessKey != "" {
		config.AwsSecretAccessKey = AwsSecretAccessKey
	}

	if AwsAccessKeyId := os.Getenv("AWS_ACCESS_KEY_ID"); AwsAccessKeyId != "" {
		config.AwsAccessKeyId = AwsAccessKeyId
	}

	if AwsRegion := os.Getenv("AWS_REGION"); AwsRegion != "" {
		config.AwsRegion = AwsRegion
	}

	return &config

}
