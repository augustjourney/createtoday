package controller

import (
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/entity"
	"createtodayapi/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	repo storage.Products
}

func (c *ProductController) GetUsersProducts(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*entity.User)
	products, err := c.repo.GetUsersProducts(context.Background(), user.ID)
	if err != nil {
		return common.DoApiResponse(ctx, 500, nil, err)
	}
	return common.DoApiResponse(ctx, 200, products, nil)
}

func NewProductController(repo storage.Products) *ProductController {
	return &ProductController{
		repo: repo,
	}
}
