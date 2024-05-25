package app

import (
	"createtodayapi/internal/config"
	"createtodayapi/internal/hero"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func New(db *sqlx.DB, redis *redis.Client, config *config.Config) *fiber.App {

	app := fiber.New()

	hero.NewHeroApp(db, redis, config, app)

	return app
}
