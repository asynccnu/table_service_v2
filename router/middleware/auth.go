package middleware

import (
	"github.com/asynccnu/table_service_v2/handler"
	"github.com/asynccnu/table_service_v2/pkg/errno"
	"github.com/asynccnu/table_service_v2/pkg/token"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := token.ParseRequest(c); err != nil {
			handler.SendUnauthorized(c, errno.ErrAuthorizationInvalid, nil)
			c.Abort()
			return
		}

		c.Next()

	}
}
