package main

import (
	"createtodayapi/internal/app"
	"createtodayapi/internal/config"
	"createtodayapi/internal/infra"
	"createtodayapi/internal/logger"
	"fmt"
	"github.com/pressly/goose/v3"
)

func main() {
	conf := config.New("")
	log := logger.New()

	db, err := infra.InitPostgres(conf.DatabaseDSN)
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = goose.SetDialect("postgres")
	if err != nil {
		log.Error(err.Error())
		return
	}

	if conf.Env == "dev" {
		version, err := goose.GetDBVersion(db.DB)
		if err != nil {
			log.Error(err.Error())
			return
		}

		log.Info(fmt.Sprintf("database version: %v", version))

		err = goose.Up(db.DB, "db/migrations")
		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	redis, err := infra.InitRedis(conf.RedisHost, conf.RedisPort)
	if err != nil {
		log.Error(err.Error())
	}

	server := app.New(db, redis, conf)

	err = server.Listen(conf.ServerAddress)

	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info(fmt.Sprintf("Server started on %s", conf.ServerAddress))
}
