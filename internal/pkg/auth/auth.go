package auth

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"github.com/duyike/greddit/internal/pkg/constant"
)

func GenerateJWT(userID int64) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["exp"] = time.Now().Add(time.Hour * 24 * time.Duration(constant.LoginExpireDays)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GetAuthenticatedUserID(c *fiber.Ctx) (uid int64, ok bool) {
	user := c.Locals("user").(*jwt.Token)
	if user == nil {
		return 0, false
	}
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))
	return userID, true
}
