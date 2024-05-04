package hero

import (
	"createtodayapi/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func NewHeroApp(db *sqlx.DB, config *config.Config, app *fiber.App) *fiber.App {

	postgres := NewPostgresRepo(db)
	memory := NewMemoryRepo()

	emailsService := NewEmailService(config, memory)
	service := NewService(postgres, config, emailsService)

	controller := NewController(service)

	hero := app.Group("/hero")

	hero.Post("/auth/login", controller.Login)
	hero.Post("/auth/login/get-magic-link", controller.GetMagicLink)
	hero.Post("/auth/login/validate-magic-link", controller.ValidateMagicLink)
	hero.Post("/auth/signup", controller.Signup)
	hero.Get("/profile", AuthMiddleware(service), controller.GetProfile)
	hero.Get("/courses", AuthMiddleware(service), controller.GetUserAccessibleProducts)
	hero.Get("/courses/:slug/lessons", AuthMiddleware(service), controller.GetUserAccessibleProduct)
	hero.Get("/courses/:courseSlug/lessons/:slug", AuthMiddleware(service), controller.GetUserAccessibleLesson)

	return app
}
