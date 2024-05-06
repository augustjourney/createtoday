package app

import (
	"createtodayapi/internal/config"
	"createtodayapi/internal/hero"
	"github.com/jmoiron/sqlx"

	"github.com/gofiber/fiber/v2"
)

func New(db *sqlx.DB, config *config.Config) *fiber.App {

	app := fiber.New()

	hero.NewHeroApp(db, config, app)

	return app
}
