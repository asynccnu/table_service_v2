package middleware

import (
	"encoding/base64"
	"github.com/asynccnu/table_service_v2/handler"
	"github.com/asynccnu/table_service_v2/pkg/errno"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		b, err := base64.StdEncoding.DecodeString(auth)
		if err != nil {
			handler.SendUnauthorized(c, errno.ErrTokenInvalid, nil, err.Error())
			c.Abort()
			return
		}

		arr := strings.Split(string(b), ":")

		if arr[0] != c.GetHeader("sid") {
			handler.SendUnauthorized(c, errno.ErrTokenInvalid, nil, "Authorization error.")
			c.Abort()
			return
		}

		c.Next()
	}
}
