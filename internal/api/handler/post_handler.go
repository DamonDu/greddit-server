package handler

import (
	"github.com/bitrise-io/go-utils/stringutil"
	"github.com/gofiber/fiber/v2"

	"github.com/duyike/greddit/internal/api/middleware"
	"github.com/duyike/greddit/internal/model"
	"github.com/duyike/greddit/internal/pkg/constant"
	"github.com/duyike/greddit/internal/service"
	"github.com/duyike/greddit/pkg/maths"
)

type PostHandler struct {
	*fiber.App
	postApp service.PostService
}

func NewPostHandler(postApp service.PostService) PostHandler {
	handler := PostHandler{
		App:     fiber.New(),
		postApp: postApp,
	}
	handler.Post("/pageQuery", handler.PageQuery)
	handler.Post("/create", middleware.UserAuth(), handler.Create)
	return handler
}

func (h *PostHandler) PageQuery(ctx *fiber.Ctx) error {
	type json struct {
		Page     int `json:"page" binding:"required"`
		PageSize int `json:"pageSize" binding:"required"`
	}
	var body json
	err := ctx.BodyParser(&body)
	if err != nil {
		panic(err)
	}
	postUsers, err := h.postApp.PageQueryPostUser(body.Page, body.PageSize+1)
	if err != nil {
		panic(err)
	}
	var realPostUsers = postUsers[:maths.Min(len(postUsers), body.PageSize)]
	err = ctx.JSON(fiber.Map{
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
	if err != nil {
		panic(err)
	}
	return nil
}

func (h *PostHandler) Create(ctx *fiber.Ctx) error {
	type json struct {
		Title string `json:"title" binding:"required"`
		Text  string `json:"text" binding:"required"`
	}
	var body json
	err := ctx.BodyParser(&body)
	if err != nil {
		panic(err)
	}
	_, err = h.postApp.Create(ctx.Locals(constant.UidHttpKey).(int64), body.Title, body.Text)
	if err != nil {
		panic(err)
	}
	return nil
}
