package utils

import (
	"github.com/damondu/greddit/domain/entity"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetUid(c *gin.Context) int64 {
	return c.GetInt64("uid")
}

func SetContextCookies(c *gin.Context, user *entity.User) {
	// directly set uid to cookies, unsafe but easy
	c.SetCookie("user_id", strconv.FormatInt(user.Uid, 10), 30*60, "/", "localhost", false, true)
}
