package handler

import (
	"github.com/damondu/greddit/application"
	"github.com/damondu/greddit/domain/entity"
	. "github.com/damondu/greddit/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	app application.UserApp
}

func NewUserHandler(apps *application.App) *UserHandler {
	return &UserHandler{app: apps.User}
}

func (h *UserHandler) Me(c *gin.Context) {
	uid := GetUid(c)
	me, meErr := h.app.Me(uid)
	if meErr != nil {
		c.Error(meErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"username": me.Username,
		"email":    me.Email,
	})
}

func (h *UserHandler) Register(c *gin.Context) {
	type json struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}
	var body json
	bindErr := c.ShouldBindJSON(&body)
	if bindErr != nil {
		c.Error(bindErr)
		return
	}
	user, appErr := h.app.Register(body.Username, body.Email, body.Password)
	if appErr != nil {
		c.Error(appErr)
		return
	}
	SetContextCookies(c, &user)
}

func (h *UserHandler) Login(c *gin.Context) {
	type json struct {
		Username string `json:"username" binding:"required_without=Email"`
		Email    string `json:"email" binding:"required_without=Username"`
		Password string `json:"password" binding:"required"`
	}
	var body json
	bindErr := c.Bind(&body)
	if bindErr != nil {
		c.Error(bindErr)
		return
	}

	var user entity.User
	var appErr error
	if body.Username != "" {
		user, appErr = h.app.LoginByUsername(body.Username, body.Password)
	} else {
		user, appErr = h.app.LoginByEmail(body.Email, body.Password)
	}
	if appErr != nil {
		c.Error(appErr)
		return
	}
	SetContextCookies(c, &user)
}
