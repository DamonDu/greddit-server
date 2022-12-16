package handler

import (
	"github.com/bitrise-io/go-utils/stringutil"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/duyike/greddit/internal/api/middleware"
	"github.com/duyike/greddit/internal/model"
	"github.com/duyike/greddit/internal/pkg/auth"
	"github.com/duyike/greddit/internal/service"
	"github.com/duyike/greddit/pkg/maths"
)

type PostHandler struct {
	*fiber.App
}

func NewPostHandler() PostHandler {
	handler := PostHandler{
		App: fiber.New(),
	}
	handler.Post("/pageQuery", handler.PageQuery)
	handler.Post("/create", middleware.UserAuth(), handler.Create)
	return handler
}

func (h *PostHandler) PageQuery(ctx *fiber.Ctx) error {
	var body struct {
		Page     int `json:"page" validate:"required"`
		PageSize int `json:"pageSize" validate:"required"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	postUsers, err := service.Post.PageQueryPostUser(body.Page, body.PageSize+1)
	if err != nil {
		return err
	}
	var realPostUsers = postUsers[:maths.Min(len(postUsers), body.PageSize)]
	return ctx.JSON(fiber.Map{
		"hasMore": len(postUsers) > body.PageSize,
		"list": realPostUsers.MapInterface(func(p *model.WithUser) interface{} {
			return fiber.Map{
				"postId":    p.PostId,
				"title":     p.Title,
				"text":      stringutil.MaxFirstCharsWithDots(p.Text, 250),
				"voteCount": p.VoteCount,
				"uid":       p.Uid,
				"username":  p.Username,
			}
		}),
	})
}

func (h *PostHandler) Create(ctx *fiber.Ctx) error {
	var body struct {
		Title string `json:"title" validate:"required"`
		Text  string `json:"text" validate:"required"`
	}
	if err := ctx.BodyParser(&body); err != nil {
		return err
	}
	if err := validator.New().Struct(body); err != nil {
		return err
	}
	uid, _ := auth.GetAuthenticatedUserID(ctx)
	if _, err := service.Post.Create(uid, body.Title, body.Text); err != nil {
		return err
	}
	return nil
}
