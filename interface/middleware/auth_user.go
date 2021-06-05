package middleware

import (
	. "github.com/damondu/greddit/domain/error"
	"github.com/gin-gonic/gin"
)

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetInt64("uid")
		if uid == 0 {
			c.Error(&ApplicationError{
				Code: NoLoginError,
				Msg:  "Please Login First",
			})
			c.Abort()
		}
		c.Next()
	}
}
