package main

import (
	"createtodayapi/internal/app"
	"createtodayapi/internal/config"
	"createtodayapi/internal/infra"
	"createtodayapi/internal/logger"
	"fmt"
)

func main() {
	conf := config.New()
	log := logger.New()

	db, err := infra.InitPostgres(conf.DatabaseDSN)
	if err != nil {
		log.Error(err.Error())
		return
	}

	s := app.New(db, conf)

	err = s.Listen(conf.ServerAddress)

	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Info(fmt.Sprintf("Server started on %s", conf.ServerAddress))
}
