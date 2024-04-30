package controller

import (
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/dto"
	"createtodayapi/internal/logger"
	"createtodayapi/internal/service"
	"encoding/json"
	"errors"
	"fmt"
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

	result, err := c.service.Login(context.Background(), &body)

	if errors.Is(err, common.ErrWrongCredentials) {
		return common.DoApiResponse(ctx, 400, nil, err)
	}

	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.DoApiResponse(ctx, 500, nil, common.ErrInternalError)
	}

	tokenCookie := new(fiber.Cookie)
	tokenCookie.Name = "token"
	tokenCookie.Value = result.Token

	ctx.Cookie(tokenCookie)

	return common.DoApiResponse(ctx, 200, result, nil)
}

func (c *AuthController) Signup(ctx *fiber.Ctx) error {
	var body dto.SignupBody

	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, 400, nil, err)
	}

	err = body.Validate()
	if err != nil {
		return common.DoApiResponse(ctx, 400, nil, err)
	}

	result, err := c.service.Signup(context.Background(), &body)
	if err != nil {
		logger.Log.Error(err.Error(), "error", err)
		return common.DoApiResponse(ctx, 500, nil, common.ErrInternalError)
	}

	if result.Token != "" {
		tokenCookie := new(fiber.Cookie)
		tokenCookie.Name = "token"
		tokenCookie.Value = result.Token

		ctx.Cookie(tokenCookie)
	}

	return common.DoApiResponse(ctx, 200, result, nil)
}

func (c *AuthController) Logout(ctx *fiber.Ctx) error {
	return nil
}

func (c *AuthController) GetMagicLink(ctx *fiber.Ctx) error {
	var body dto.GetMagicLinkBody
	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, 400, nil, err)
	}

	err = body.Validate()
	if err != nil {
		return common.DoApiResponse(ctx, 400, nil, err)
	}

	err = c.service.GetMagicLink(context.Background(), body.Email)

	if err != nil {
		return common.DoApiResponse(ctx, 500, nil, err)
	}

	type GetMagicLinkResult struct {
		Message string `json:"message"`
	}

	return common.DoApiResponse(ctx, 200, GetMagicLinkResult{
		Message: "Письмо с ссылкой для входа без пароля отправлено на вашу почту",
	}, nil)
}

func (c *AuthController) ValidateMagicLink(ctx *fiber.Ctx) error {
	var body dto.ValidateMagicLinkBody
	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, 400, nil, err)
	}

	err = body.Validate()
	if err != nil {
		return common.DoApiResponse(ctx, 400, nil, err)
	}

	result, err := c.service.ValidateMagicLink(context.Background(), body.Token)

	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, common.ErrMagicLinkExpired) || errors.Is(err, common.ErrInvalidMagicLink) {
			return common.DoApiResponse(ctx, 400, nil, err)
		}
		return common.DoApiResponse(ctx, 500, nil, err)
	}

	tokenCookie := new(fiber.Cookie)
	tokenCookie.Name = "token"
	tokenCookie.Value = result.Token

	ctx.Cookie(tokenCookie)

	return common.DoApiResponse(ctx, 200, result, nil)
}

func NewAuthController(service *service.Auth) *AuthController {
	return &AuthController{
		service: service,
	}
}
