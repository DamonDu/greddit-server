package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

func UserAuth() fiber.Handler {
	return userAuth(false)
}

func WeakUserAuth() fiber.Handler {
	return userAuth(true)
}

func userAuth(weakCheck bool) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if weakCheck {
				return c.JSON(nil)
			}
			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusUnauthorized).
					JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
			}
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
		},
	})
}
