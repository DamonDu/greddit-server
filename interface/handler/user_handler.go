package handler

import (
	"errors"
	"github.com/damondu/greddit/application"
	"github.com/damondu/greddit/domain/entity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type UserHandler struct {
	app *application.UserApp
}

func NewUserHandler(apps *application.App) *UserHandler {
	return &UserHandler{app: apps.User}
}

func (h *UserHandler) Me(c *gin.Context) {
	uid := c.GetInt64("uid")
	me, meErr := h.app.Me(uid)
	if meErr != nil {
		if errors.Is(meErr, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, nil)
		} else {
			c.Error(meErr)
		}
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
	bindErr := c.Bind(&body)
	if bindErr != nil {
		c.Error(bindErr)
		return
	}
	user, appErr := h.app.Register(body.Username, body.Email, body.Password)
	if appErr != nil {
		c.Error(appErr)
		return
	}
	cookie := entity.NewUserCookie(&user)
	c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.Secure)
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
	cookie := entity.NewUserCookie(&user)
	c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.Secure)
}

func (h *UserHandler) Logout(c *gin.Context) {
	cookie := entity.NewUserCookie(nil)
	c.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
}
