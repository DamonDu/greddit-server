package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/duyike/greddit/internal/api/middleware"
	"github.com/duyike/greddit/internal/model"
	"github.com/duyike/greddit/internal/pkg/auth"
	"github.com/duyike/greddit/internal/service"
)

type UserHandler struct {
	*fiber.App
	userApp service.UserService
}

func NewUserHandler(userApp service.UserService) UserHandler {
	handler := UserHandler{
		App:     fiber.New(),
		userApp: userApp,
	}
	handler.Post("/me", middleware.WeakUserAuth(), handler.Me)
	handler.Post("/register", handler.Register)
	handler.Post("/login", handler.Login)
	handler.Post("/logout", handler.Logout)
	return handler
}

func (h *UserHandler) Me(ctx *fiber.Ctx) error {
	uid := auth.GetAuthenticatedUserID(ctx)
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
	token, err := auth.GenerateJWT(registerUser.Uid)
	if err != nil {
		panic(err)
	}
	return ctx.JSON(fiber.Map{
		"token": token,
	})
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

	var u model.User
	var appErr error
	if body.Username != "" {
		u, appErr = h.userApp.LoginByUsername(body.Username, body.Password)
	} else {
		u, appErr = h.userApp.LoginByEmail(body.Email, body.Password)
	}
	if appErr != nil {
		panic(appErr)
	}
	token, err := auth.GenerateJWT(u.Uid)
	if err != nil {
		panic(err)
	}
	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func (h *UserHandler) Logout(ctx *fiber.Ctx) error {
	return nil
}
