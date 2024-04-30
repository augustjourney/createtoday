package middleware

import (
	"context"
	"createtodayapi/internal/common"
	"createtodayapi/internal/logger"
	"createtodayapi/internal/service"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func Auth(ctx *fiber.Ctx, auth *service.Auth) error {
	authHeader := ctx.Get("authorization")

	if authHeader == "" {
		return common.DoApiResponse(ctx, 401, nil, nil)
	}

	authHeaderData := strings.SplitAfterN(authHeader, "Bearer ", 2)

	if len(authHeaderData) < 2 {
		logger.Log.Warn(fmt.Sprintf("Could not split Bearer token %s", authHeader))
		return common.DoApiResponse(ctx, 401, nil, nil)
	}

	token := authHeaderData[1]

	user, err := auth.ValidateJWTToken(context.Background(), token)

	if err != nil {
		if errors.Is(err, common.ErrTokenExpired) {
			logger.Log.Warn("Token expired")
			return common.DoApiResponse(ctx, 403, nil, nil)
		}

		if errors.Is(err, common.ErrInvalidToken) {
			logger.Log.Warn("Invalid token")
			return common.DoApiResponse(ctx, 403, nil, nil)
		}

		logger.Log.Error(fmt.Sprintf("Not valid token %s", err.Error()))
		return common.DoApiResponse(ctx, 500, nil, err)
	}

	if user == nil {
		logger.Log.Warn("Not found user by token", "token", token)
		return common.DoApiResponse(ctx, 401, nil, err)
	}

	ctx.Locals("user", user)

	return ctx.Next()
}
