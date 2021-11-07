package handler

import (
	"github.com/bitrise-io/go-utils/stringutil"
	"github.com/gofiber/fiber/v2"

	"github.com/duyike/greddit/internal/api/middleware"
	"github.com/duyike/greddit/internal/post"
	math2 "github.com/duyike/greddit/pkg/math"
)

type PostHandler struct {
	*fiber.App
	postApp post.App
}

func NewPostHandler(postApp post.App) PostHandler {
	handler := PostHandler{
		App:     fiber.New(),
		postApp: postApp,
	}
	handler.Post("/pageQuery", handler.PageQuery)
	handler.Post("/create", middleware.NewSimpleAuth(true), handler.Create)
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
	var realPostUsers = postUsers[:math2.Min(len(postUsers), body.PageSize)]
	err = ctx.JSON(fiber.Map{
		"hasMore": len(postUsers) > body.PageSize,
		"list": realPostUsers.MapInterface(func(p *post.WithUser) interface{} {
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
	_, err = h.postApp.Create(ctx.Locals(uidHttpKey).(int64), body.Title, body.Text)
	if err != nil {
		panic(err)
	}
	return nil
}
