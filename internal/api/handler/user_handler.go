package handler

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/duyike/greddit/internal/api/middleware"
	"github.com/duyike/greddit/internal/user"
)

var (
	uidCookieName = os.Getenv("UID_COOKIE_NAME")
	uidHttpKey    = os.Getenv("UID_HTTP_KEY")
)

type UserHandler struct {
	*fiber.App
	userApp user.App
}

func NewUserHandler(userApp user.App) UserHandler {
	handler := UserHandler{
		App:     fiber.New(),
		userApp: userApp,
	}
	handler.Post("/me", middleware.NewSimpleAuth(false), handler.Me)
	handler.Post("/register", handler.Register)
	handler.Post("/login", handler.Login)
	handler.Post("/logout", handler.Logout)
	return handler
}

func (h *UserHandler) Me(ctx *fiber.Ctx) error {
	uid := middleware.GetUid(ctx)
	if uid == nil {
		return ctx.JSON(nil)
	}
	me, err := h.userApp.QueryByUid(*uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		} else {
			panic(err)
		}
	}
	return ctx.JSON(fiber.Map{
		"username": me.Username,
		"email":    me.Email,
	})
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	type json struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}
	var body json
	err := ctx.BodyParser(&body)
	if err != nil {
		panic(err)
	}
	registerUser, err := h.userApp.Register(body.Username, body.Email, body.Password)
	if err != nil {
		panic(err)
	}
	ctx.Cookie(h.buildCookie(&registerUser))
	return nil
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	type json struct {
		Username string `json:"username" binding:"required_without=Email"`
		Email    string `json:"email" binding:"required_without=Username"`
		Password string `json:"password" binding:"required"`
	}
	var body json
	err := ctx.BodyParser(&body)
	if err != nil {
		panic(err)
	}

	var u user.User
	var appErr error
	if body.Username != "" {
		u, appErr = h.userApp.LoginByUsername(body.Username, body.Password)
	} else {
		u, appErr = h.userApp.LoginByEmail(body.Email, body.Password)
	}
	if appErr != nil {
		panic(appErr)
	}
	ctx.Cookie(h.buildCookie(&u))
	return nil
}

func (h *UserHandler) Logout(ctx *fiber.Ctx) error {
	ctx.Cookie(h.buildCookie(nil))
	return nil
}

func (h UserHandler) buildCookie(user *user.User) *fiber.Cookie {
	var (
		value  = ""
		maxAge = -1
	)
	if user != nil {
		value = strconv.FormatInt(user.Uid, 10)
		maxAge = 1800
	}
	return &fiber.Cookie{
		Name:     uidCookieName,
		Value:    value,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   maxAge,
		Expires:  time.Time{},
		Secure:   false,
		HTTPOnly: true,
	}
}
