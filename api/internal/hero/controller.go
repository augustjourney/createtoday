package hero

import (
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/logger"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type IController interface {
	// Auth
	Login(ctx *fiber.Ctx) error
	Signup(ctx *fiber.Ctx) error
	GetMagicLink(ctx *fiber.Ctx) error
	ValidateMagicLink(ctx *fiber.Ctx) error

	// Products
	GetUserAccessibleProducts(ctx *fiber.Ctx) error

	// Profile
	GetProfile(ctx *fiber.Ctx) error
}

type Controller struct {
	service IService
}

func (c *Controller) Login(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "application/json")
	var body LoginBody

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

func (c *Controller) Signup(ctx *fiber.Ctx) error {
	var body SignupBody

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

	if result.Token != nil {
		tokenCookie := new(fiber.Cookie)
		tokenCookie.Name = "token"
		tokenCookie.Value = *result.Token

		ctx.Cookie(tokenCookie)
	}

	return common.DoApiResponse(ctx, 200, result, nil)
}

func (c *Controller) GetMagicLink(ctx *fiber.Ctx) error {
	var body GetMagicLinkBody
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

func (c *Controller) ValidateMagicLink(ctx *fiber.Ctx) error {
	var body ValidateMagicLinkBody
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

func (c *Controller) GetProfile(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*User)
	profile, err := c.service.GetProfile(context.Background(), user.ID)

	if err != nil {
		return common.DoApiResponse(ctx, 500, nil, err)
	}

	if profile == nil {
		return common.DoApiResponse(ctx, 404, nil, common.ErrUserNotFound)
	}

	return common.DoApiResponse(ctx, 200, profile, nil)
}

func (c *Controller) UpdateProfile(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*User)
	var body UpdateProfileBody
	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil {
		logger.Log.Error(err.Error())
		return common.DoApiResponse(ctx, 400, nil, err)
	}

	// TODO: валидация body
	// Все поля из структуры должны быть в json-body
	// Например: если в json-body нет полей,
	// Которые должны быть в структуре
	// А в БД эти поля заполнены
	// Получается перезатрем их на null
	// При этом если поле есть в json-body со значением null
	// Это корректно

	err = c.service.UpdateProfile(context.Background(), user.ID, body)

	if err != nil {
		return common.DoApiResponse(ctx, 500, nil, err)
	}

	return common.DoApiResponse(ctx, 200, "Профиль обновлен", nil)
}

func (c *Controller) GetUserAccessibleProducts(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*User)
	products, err := c.service.GetUserAccessibleProducts(context.Background(), user.ID)
	if err != nil {
		return common.DoApiResponse(ctx, 500, nil, err)
	}
	return common.DoApiResponse(ctx, 200, products, nil)
}

func (c *Controller) GetUserAccessibleProduct(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*User)
	slug := ctx.Params("slug")
	product, err := c.service.GetUserAccessibleProduct(context.Background(), slug, user.ID)
	if errors.Is(err, common.ErrProductNotFound) {
		return common.DoApiResponse(ctx, 404, nil, common.ErrProductNotFound)
	}
	if err != nil {
		return common.DoApiResponse(ctx, 500, nil, err)
	}
	return common.DoApiResponse(ctx, 200, product, nil)
}

func (c *Controller) GetUserAccessibleLesson(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*User)
	slug := ctx.Params("slug")
	lesson, err := c.service.GetUserAccessibleLesson(context.Background(), slug, user.ID)
	if errors.Is(err, common.ErrLessonNotFound) {
		return common.DoApiResponse(ctx, 404, nil, err)
	}
	if err != nil {
		return common.DoApiResponse(ctx, 500, nil, err)
	}
	return common.DoApiResponse(ctx, 200, lesson, nil)
}

func NewController(service IService) *Controller {
	return &Controller{
		service: service,
	}
}
