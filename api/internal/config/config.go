package config

import (
	"createtodayapi/internal/logger"
	"flag"
	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	BaseURL             string
	ServerAddress       string
	DatabaseDSN         string `env:"DATABASE_DSN"`
	JwtTokenExp         time.Duration
	MagicLinkExp        time.Duration
	JwtTokenSecretKey   string `env:"JWT_TOKEN_SECRET_KEY"`
	JwtSigningMethod    jwt.SigningMethod
	HeroAppBaseURL      string `env:"HERO_APP_BASE_URL"`
	AwsSecretAccessKey  string `env:"AWS_SECRET_ACCESS_KEY"`
	AwsAccessKeyId      string `env:"AWS_ACCESS_KEY_ID"`
	AwsRegion           string `env:"AWS_REGION"`
	Env                 string `env:"ENV"` // dev, stage, prod
	S3Endpoint          string
	S3Region            string
	S3AccessKeyId       string `env:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey   string `env:"S3_SECRET_ACCESS_KEY"`
	CdnUrl              string `env:"CDN_URL"`
	PhotosBucket        string
	VideosBucket        string
	S3Provider          string
	TinkoffTestLogin    string `env:"TINKOFF_TEST_LOGIN"`
	TinkoffTestPassword string `env:"TINKOFF_TEST_PASSWORD"`
}

var config *Config

func (c *Config) loadEnv(path string) {
	if path == "" {
		path = ".env"
	}
	err := godotenv.Load(path)
	if err != nil {
		logger.Log.Error("No .env file found", "error", err.Error())
	}
}

func (c *Config) parseEnv() {
	var flagServerAddress = flag.String("a", "localhost:8080", "Server address on which server is running")
	var flagDatabaseDSN = flag.String("d", "", "Database DSN")
	var flagEnv = flag.String("e", "", "Environment")

	c.JwtSigningMethod = jwt.SigningMethodHS256
	c.JwtTokenExp = time.Hour * 720
	c.MagicLinkExp = time.Minute * 1
	c.ServerAddress = *flagServerAddress
	c.Env = "dev"
	c.S3Endpoint = "https://s3.storage.selcloud.ru"
	c.S3Region = "ru-1a"
	c.PhotosBucket = "photos"
	c.VideosBucket = "videos"
	c.S3Provider = "selectel"

	err := env.Parse(c)
	if err != nil {
		logger.Log.Error("could not parse env vars:", "error", err)
	}

	if *flagDatabaseDSN != "" {
		c.DatabaseDSN = *flagDatabaseDSN
	}

	if *flagEnv != "" {
		c.Env = *flagEnv
	}
}

func New(pathToEnv string) *Config {
	if config != nil {
		return config
	}

	config = &Config{}

	config.loadEnv(pathToEnv)
	config.parseEnv()

	return config

}
