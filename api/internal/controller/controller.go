package controller

import "github.com/gofiber/fiber/v2"

type Auth interface {
	Login(ctx *fiber.Ctx) error
	Signup(ctx *fiber.Ctx) error
	GetMagicLink(ctx *fiber.Ctx) error
	ValidateMagicLink(ctx *fiber.Ctx) error
}

type Products interface {
	GetUsersProducts(ctx *fiber.Ctx) error
}

type Profile interface {
	GetProfile(ctx *fiber.Ctx) error
}
