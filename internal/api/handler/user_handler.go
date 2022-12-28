package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/duyike/greddit/internal/model"
	"github.com/duyike/greddit/internal/pkg/api"
	"github.com/duyike/greddit/internal/service"
	bizError "github.com/duyike/greddit/pkg/errors"
)

type UserHandler struct {
	*fiber.App
}

func NewUserHandler() UserHandler {
	handler := UserHandler{
		App: fiber.New(),
	}
	handler.Post("/me", handler.Me)
	handler.Post("/register", handler.Register)
	handler.Post("/login", handler.Login)
	handler.Post("/logout", handler.Logout)
	return handler
}

func (h *UserHandler) Me(ctx *fiber.Ctx) error {
	uid, err := api.GetAuthUserID(ctx)
	if err != nil {
		return ctx.JSON(nil)
	}
	me, err := service.User.QueryByUid(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return bizError.UserNotExistsError
		} else {
			return err
		}
	}
	return ctx.JSON(fiber.Map{
		"username": me.Username,
		"email":    me.Email,
	})
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	var body struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		Email    string `json:"email" validate:"required"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	registerUser, err := service.User.Register(body.Username, body.Email, body.Password)
	if err != nil {
		return err
	}
	token, err := api.GenerateJWT(registerUser.Uid)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	var body struct {
		Username string `json:"username" validate:"required_without=Email"`
		Email    string `json:"email" validate:"required_without=Username"`
		Password string `json:"password" validate:"required"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}

	var u model.User
	var err error
	if body.Username != "" {
		u, err = service.User.LoginByUsername(body.Username, body.Password)
	} else {
		u, err = service.User.LoginByEmail(body.Email, body.Password)
	}
	if err != nil {
		return err
	}
	token, err := api.GenerateJWT(u.Uid)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

func (h *UserHandler) Logout(ctx *fiber.Ctx) error {
	return nil
}
