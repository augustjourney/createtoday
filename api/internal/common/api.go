package common

import (
	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	OK      bool        `json:"ok"`
	Result  interface{} `json:"result"`
	Message string      `json:"message,omitempty"`
}

func DoApiResponse(ctx *fiber.Ctx, status int, data interface{}, err error) error {
	var message string
	if err != nil {
		message = err.Error()
	}
	return ctx.Status(status).JSON(APIResponse{
		OK:      err == nil,
		Result:  data,
		Message: message,
	})
}
