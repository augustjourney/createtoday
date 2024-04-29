package controller

import (
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/dto"
	"createtodayapi/internal/logger"
	"createtodayapi/internal/service"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	service *service.Auth
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "application/json")
	var body dto.LoginBody

	err := json.Unmarshal(ctx.Body(), &body)

	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, 400, nil, err)
	}

	err = body.Validate()

	if err != nil {
		return common.DoApiResponse(ctx, 400, nil, err)
	}

	result, err := c.service.Login(context.Background(), body.Email, body.Password)

	if errors.As(err, &common.ErrWrongCredentials) {
		return common.DoApiResponse(ctx, 400, nil, err)
	}

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.DoApiResponse(ctx, 500, nil, common.ErrInternalError)
	}

	loginResult := dto.LoginResult{
		Token: result.Token,
	}

	tokenCookie := new(fiber.Cookie)
	tokenCookie.Name = "token"
	tokenCookie.Value = result.Token

	ctx.Cookie(tokenCookie)

	return common.DoApiResponse(ctx, 200, loginResult, nil)
}

func (c *AuthController) Signup(ctx *fiber.Ctx) error {
	return nil
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	return nil
}

func (c *AuthController) GetMagicLink(ctx *fiber.Ctx) error {
	return nil
}

func (c *AuthController) ValidateMagicLink(ctx *fiber.Ctx) error {
	return nil
}

func NewAuthController(service *service.Auth) *AuthController {
	return &AuthController{
		service: service,
	}
}