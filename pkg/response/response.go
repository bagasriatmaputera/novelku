package response

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type ErrorBody struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type Meta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type Body struct {
	Data  any        `json:"data"`
	Error *ErrorBody `json:"error"`
	Meta  *Meta      `json:"meta"`
}

func Success(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(Body{
		Data: data,
	})
}

func Created(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Body{
		Data: data,
	})
}

func List(c *fiber.Ctx, data any, meta Meta) error {
	return c.Status(fiber.StatusOK).JSON(
		Body{
			Data: data,
			Meta: &meta,
		},
	)
}

func Error(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(Body{
		Error: &ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := "Something Went Wrong!"

	var fe *fiber.Error
	if errors.As(err, &fe) {
		status = fe.Code
		message = fe.Message
	}

	return Error(c, status, "ERROR", message)
}

func ValidatationError(c *fiber.Ctx, details map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(Body{
		Error: &ErrorBody{
			Code:    "VALIDATION ERROR",
			Message: "Data yang dikirim tidak valid",
			Details: details,
		},
	})
}
