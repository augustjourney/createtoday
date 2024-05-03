package app

import (
	"createtodayapi/internal/config"
	"createtodayapi/internal/controller"
	"createtodayapi/internal/middleware"
	"createtodayapi/internal/service"
	"createtodayapi/internal/storage/memory"
	"createtodayapi/internal/storage/postgres"

	"github.com/jmoiron/sqlx"

	"github.com/gofiber/fiber/v2"
)

func New(db *sqlx.DB, config *config.Config) *fiber.App {

	// создать все репозитории
	usersRepo := postgres.NewUsersRepo(db)
	emailsRepo := memory.NewEmailsRepo()
	productsRepo := postgres.NewProductsRepo(db)

	// создать все сервисы
	emailService := service.NewEmailService(config, emailsRepo)
	profileService := service.NewProfileService(usersRepo)
	authService := service.NewAuthService(usersRepo, config, emailService)

	// создать все контроллеры
	profileController := controller.NewProfileController(profileService)
	authController := controller.NewAuthController(authService)
	productsController := controller.NewProductController(productsRepo)

	app := fiber.New()

	app.Post("/hero/auth/login", authController.Login)
	app.Post("/hero/auth/login/get-magic-link", authController.GetMagicLink)
	app.Post("/hero/auth/login/validate-magic-link", authController.ValidateMagicLink)
	app.Post("/hero/auth/signup", authController.Signup)
	app.Get("/hero/profile", middleware.Auth(authService), profileController.GetProfile)
	app.Get("/hero/courses", middleware.Auth(authService), productsController.GetUsersProducts)

	return app
}
