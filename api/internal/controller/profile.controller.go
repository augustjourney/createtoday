package controller

import (
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/entity"
	"createtodayapi/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ProfileController struct {
	service *service.Profile
}

func (c *ProfileController) GetProfile(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*entity.User)
	profile, err := c.service.GetProfile(context.Background(), user.ID)

	if err != nil {
		return common.DoApiResponse(ctx, 500, nil, err)
	}

	if profile == nil {
		return common.DoApiResponse(ctx, 404, nil, common.ErrUserNotFound)
	}

	return common.DoApiResponse(ctx, 200, profile, nil)
}

func NewProfileController(service *service.Profile) *ProfileController {
	return &ProfileController{
		service: service,
	}
}
