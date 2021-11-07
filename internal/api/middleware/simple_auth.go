package middleware

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"

	error2 "github.com/duyike/greddit/pkg/errors"
)

var (
	uidCookieName = os.Getenv("UID_COOKIE_NAME")
	uidHttpKey    = os.Getenv("UID_HTTP_KEY")
)

func NewSimpleAuth(strongAuth bool) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		parseCookies(ctx)
		if strongAuth {
			CheckLogin(ctx)
		}
		return ctx.Next()
	}
}

func CheckLogin(ctx *fiber.Ctx) int64 {
	uid := GetUid(ctx)
	if uid == nil {
		panic(error2.NoLoginError)
	} else {
		return *uid
	}
}

func GetUid(ctx *fiber.Ctx) *int64 {
	if uid, ok := ctx.Locals(uidHttpKey).(int64); !ok {
		return nil
	} else {
		return &uid
	}
}

func parseCookies(ctx *fiber.Ctx) {
	uidString := ctx.Cookies(uidCookieName, "")
	if uidString == "" {
		return
	}
	uid, err := strconv.ParseInt(uidString, 10, 64)
	if err != nil {
		return
	}
	ctx.Locals(uidHttpKey, uid)
}
