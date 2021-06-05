package handler

import (
	"github.com/bitrise-io/go-utils/stringutil"
	"github.com/damondu/greddit/application"
	"github.com/damondu/greddit/domain/entity"
	"github.com/damondu/greddit/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostHandler struct {
	app *application.PostApp
}

func NewPostHandler(app *application.App) *PostHandler {
	return &PostHandler{app: app.Post}
}

func (h *PostHandler) PageQuery(c *gin.Context) {
	type json struct {
		Page     int `json:"page" binding:"required"`
		PageSize int `json:"pageSize" binding:"required"`
	}
	var body json
	bindErr := c.Bind(&body)
	if bindErr != nil {
		c.Error(bindErr)
		return
	}

	postUsers, appErr := h.app.PageQueryPostUser(body.Page, body.PageSize+1)
	if appErr != nil {
		c.Error(appErr)
		return
	}
	var realPostUsers = postUsers[:utils.Min(len(postUsers), body.PageSize)]
	c.JSON(http.StatusOK, gin.H{
		"hasMore": len(postUsers) > body.PageSize,
		"list": realPostUsers.MapInterface(func(p *entity.PostUser) interface{} {
			return gin.H{
				"postId":    p.Post.PostId,
				"title":     p.Post.Title,
				"text":      stringutil.MaxLastCharsWithDots(p.Post.Text, 250),
				"voteCount": p.Post.VoteCount,
				"uid":       p.User.Uid,
				"username":  p.User.Username,
			}
		}),
	})

}

func (h *PostHandler) Create(c *gin.Context) {
	type json struct {
		Title string `json:"title" binding:"required"`
		Text  string `json:"text" binding:"required"`
	}
	var body json
	bindErr := c.Bind(&body)
	if bindErr != nil {
		c.Error(bindErr)
		return
	}
	_, appErr := h.app.Create(c.GetInt64("uid"), body.Title, body.Text)
	if appErr != nil {
		c.Error(appErr)
		return
	}
}
