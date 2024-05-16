package config

import (
	"createtodayapi/internal/logger"
	"flag"
	"github.com/caarlos0/env/v9"
	"time"

	"github.com/joho/godotenv"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	BaseURL             string        `env:"BASE_URL"`
	ServerAddress       string        `env:"SERVER_ADDRESS"`
	DatabaseDSN         string        `env:"DATABASE_DSN"`
	JwtTokenExp         time.Duration `env:"JWT_TOKEN_EXP"`
	MagicLinkExp        time.Duration `env:"MAGIC_LINK_EXP"`
	JwtTokenSecretKey   string        `env:"JWT_TOKEN_SECRET_KEY"`
	JwtSigningMethod    jwt.SigningMethod
	HeroAppBaseURL      string `env:"HERO_APP_BASE_URL"`
	AwsSecretAccessKey  string `env:"AWS_SECRET_ACCESS_KEY"`
	AwsAccessKeyId      string `env:"AWS_ACCESS_KEY_ID"`
	AwsRegion           string `env:"AWS_REGION"`
	Env                 string `env:"ENV"` // dev, stage, prod
	S3Endpoint          string `env:"S3_ENDPOINT"`
	S3Region            string `env:"S3_REGION"`
	S3AccessKeyId       string `env:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey   string `env:"S3_SECRET_ACCESS_KEY"`
	CdnUrl              string `env:"CDN_URL"`
	PhotosBucket        string `env:"PHOTOS_BUCKET"`
	VideosBucket        string `env:"VIDEOS_BUCKET"`
	S3Provider          string `env:"S3_PROVIDER"`
	TinkoffTestLogin    string `env:"TINKOFF_TEST_LOGIN"`
	TinkoffTestPassword string `env:"TINKOFF_TEST_PASSWORD"`
}

var config Config

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Log.Error("No .env file found", "error", err.Error())
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
	config.VideosBucket = "videos"
	config.S3Provider = "selectel"

	err := env.Parse(&config)
	if err != nil {
		logger.Log.Error("could not parse env vars:", "error", err)
	}

	if *flagDatabaseDSN != "" {
		config.DatabaseDSN = *flagDatabaseDSN
	}

	if *flagEnv != "" {
		config.Env = *flagEnv
	}

	return &config

}
