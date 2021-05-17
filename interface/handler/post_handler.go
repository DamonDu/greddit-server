package handler

import (
	"github.com/damondu/greddit/application"
	"github.com/damondu/greddit/domain/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostHandler struct {
	app application.PostApp
}

func NewPostHandler(app *application.App) *PostHandler {
	return &PostHandler{app: app.Post}
}

func (h PostHandler) PageQuery(c *gin.Context) {
	type json struct {
		Page     int `json:"page" binding:"required"`
		PageSize int `json:"pageSize" binding:"required"`
	}
	var body json
	bindErr := c.ShouldBindJSON(&body)
	if bindErr != nil {
		c.Error(bindErr)
		return
	}

	postUsers, appErr := h.app.PageQueryPostUser(body.Page, body.PageSize+1)
	if appErr != nil {
		c.Error(appErr)
		return
	}
	var realPostUsers = postUsers[:body.PageSize]
	c.JSON(http.StatusOK, gin.H{
		"hasMore": len(postUsers) > body.PageSize,
		"list": realPostUsers.MapInterface(func(p *entity.PostUser) interface{} {
			return gin.H{
				"postId":    p.Post.PostId,
				"title":     p.Post.Title,
				"text":      p.Post.Text[0:250] + "...",
				"voteCount": p.Post.VoteCount,
				"uid":       p.User.Uid,
				"username":  p.User.Username,
			}
		}),
	})
}
