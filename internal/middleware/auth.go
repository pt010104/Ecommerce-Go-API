package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/jwt"
	"github.com/pt010104/api-golang/pkg/response"
)

func (m Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("x-cliend-id")

		ctx := c.Request.Context()
		k, err := m.repo.DetailKeyToken(ctx, userID)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		keyString := k.SecretKey

		tokenString := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		if tokenString == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		payload, err := jwt.Verify(tokenString, keyString)
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		ctx = jwt.SetPayloadToContext(ctx, payload)

		scope := jwt.NewScope(payload)
		ctx = jwt.SetScopeToContext(ctx, scope)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
