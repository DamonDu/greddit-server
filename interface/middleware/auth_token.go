package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		uidString, cookieErr := c.Cookie("user_id")
		if cookieErr != nil {
			c.Next()
			return
		}
		uid, parseErr := strconv.ParseInt(uidString, 10, 64)
		if parseErr != nil {
			c.Next()
			return
		}
		c.Set("uid", uid)
		c.Next()
	}
}
