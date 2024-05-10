package config

import (
	"createtodayapi/internal/logger"
	"flag"
	"os"
	"time"

	"github.com/joho/godotenv"

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
	Env                string `env:"ENV"` // dev, stage, prod
	S3Endpoint         string `env:"S3_ENDPOINT"`
	S3Region           string `env:"S3_REGION"`
	S3AccessKeyId      string `env:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey  string `env:"S3_SECRET_ACCESS_KEY"`
	CdnUrl             string `env:"CDN_URL"`
	PhotosBucket       string `env:"PHOTOS_BUCKET"`
	S3Provider         string `env:"S3_PROVIDER"`
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
	var flagEnv = flag.String("e", "", "Environment")

	config.JwtSigningMethod = jwt.SigningMethodHS256
	config.JwtTokenExp = time.Hour * 720
	config.MagicLinkExp = time.Minute * 1
	config.ServerAddress = *flagServerAddress
	config.Env = "dev"
	config.S3Endpoint = "https://s3.storage.selcloud.ru"
	config.S3Region = "ru-1a"
	config.PhotosBucket = "photos"
	config.S3Provider = "selectel"

	if databaseDSN := os.Getenv("DATABASE_DSN"); databaseDSN != "" {
		config.DatabaseDSN = databaseDSN
	}

	if *flagDatabaseDSN != "" {
		config.DatabaseDSN = *flagDatabaseDSN
	}

	if env := os.Getenv("ENV"); env != "" {
		config.Env = env
	}

	if *flagEnv != "" {
		config.Env = *flagEnv
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

	if S3AccessKeyId := os.Getenv("S3_ACCESS_KEY_ID"); S3AccessKeyId != "" {
		config.S3AccessKeyId = S3AccessKeyId
	}

	if S3SecretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY"); S3SecretAccessKey != "" {
		config.S3SecretAccessKey = S3SecretAccessKey
	}

	if CdnUrl := os.Getenv("CDN_URL"); CdnUrl != "" {
		config.CdnUrl = CdnUrl
	}

	if PhotosBucket := os.Getenv("PHOTOS_BUCKET"); PhotosBucket != "" {
		config.PhotosBucket = PhotosBucket
	}

	return &config

}
