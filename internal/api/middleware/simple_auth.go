package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/duyike/greddit/internal/pkg/constant"
	error2 "github.com/duyike/greddit/pkg/errors"
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
	if uid, ok := ctx.Locals(constant.UidHttpKey).(int64); !ok {
		return nil
	} else {
		return &uid
	}
}

func parseCookies(ctx *fiber.Ctx) {
	uidString := ctx.Cookies(constant.UidCookieName, "")
	if uidString == "" {
		return
	}
	uid, err := strconv.ParseInt(uidString, 10, 64)
	if err != nil {
		return
	}
	ctx.Locals(constant.UidHttpKey, uid)
}
