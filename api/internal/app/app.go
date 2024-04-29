package app

import (
	"createtodayapi/internal/config"
	"createtodayapi/internal/controller"
	"createtodayapi/internal/logger"
	"createtodayapi/internal/middleware"
	"createtodayapi/internal/service"
	"createtodayapi/internal/storage/postgres"

	"github.com/jmoiron/sqlx"

	"github.com/gofiber/fiber/v2"
)

func New(db *sqlx.DB) *fiber.App {
	config := config.New()
	logger := logger.New()
	// создать все репозитории
	usersRepo := postgres.NewUsersRepo(db)

	// создать все сервисы
	profileService := service.NewProfileService(usersRepo)
	authService := service.NewAuthService(usersRepo, config)

	// создать все контроллеры
	profileController := controller.NewProfileController(profileService)
	authController := controller.NewAuthController(authService)

	router := fiber.New()

	router.Get("/hero/profile", func(ctx *fiber.Ctx) error {
		return middleware.Auth(ctx, authService)
	}, profileController.GetProfile)
	router.Post("/hero/auth/login", authController.Login)

	logger.Info("App compiled")

	return router
}
