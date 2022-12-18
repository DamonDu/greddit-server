package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"

	"github.com/duyike/greddit/internal/pkg/constant"
	"github.com/duyike/greddit/pkg/errors"
)

func UserAuth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// save error to context and pass to next
			var bizErr error
			if err.Error() == "Missing or malformed JWT" {
				bizErr = errors.MissingOrMalformedJwtError
			} else {
				bizErr = errors.InvalidOrExpiredJwtError
			}
			c.Locals(constant.AuthErrorKey, bizErr)
			return c.Next()
		},
	})
}
