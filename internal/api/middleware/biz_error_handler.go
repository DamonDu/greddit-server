package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/duyike/greddit/pkg/errors"
)

func NewBizErrorHandler() func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		if e, ok := err.(*errors.BizError); ok {
			return c.Status(e.Status).JSON(e)
		}
		return fiber.DefaultErrorHandler(c, err)
	}
}
