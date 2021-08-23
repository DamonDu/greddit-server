package middleware

import (
	error2 "github.com/damondu/greddit/pkg/error"
	"github.com/gofiber/fiber/v2"
)

func NewBizErrorHandler() func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		if e, ok := err.(*error2.BizError); ok {
			return c.Status(e.Status).JSON(e)
		}
		return fiber.DefaultErrorHandler(c, err)
	}
}
