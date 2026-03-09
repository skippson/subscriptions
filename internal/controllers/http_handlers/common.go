package httphandlers

import (
	"subscriptions/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type Middleware interface {
	SetRequestID() fiber.Handler
}

// @Description Тело ответа при успехе
type SuccessResponse struct{
	Data any `json:"data"`
} // @SuccessResponse

// @Description Тело ответа при ошибке
type ErrorResponse struct {
	// Описание ошибки
	Msg    string `json:"message"`

	// Статус код
	Status int    `json:"status"`
} // @name ErrorResponse

func writeSuccess(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(SuccessResponse{
		Data: data,
	})
}

func writeError(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(ErrorResponse{
		Status: status,
		Msg:    msg,
	})
}

func getLogger(c *fiber.Ctx) logger.Logger {
	return c.Locals("logger").(logger.Logger)
}
