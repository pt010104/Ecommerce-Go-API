package middleware

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/jwt"
	"github.com/pt010104/api-golang/pkg/response"
)

func (m Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("x-client-id")
		sessionID := c.GetHeader("session-id")

		ctx := c.Request.Context()

		var wg sync.WaitGroup
		var wgErr error
		var k models.KeyToken
		var u models.User

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			k, err = m.userUC.DetailKeyToken(ctx, userID, sessionID)
			if err != nil {
				wgErr = fmt.Errorf("middleware.Auth.user.usecase.DetailKeyToken: %v", err)
				return
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			u, err = m.userUC.GetModel(ctx, userID)
			if err != nil {
				wgErr = fmt.Errorf("middleware.Auth.user.usecase.DetailUser: %v", err)
				return
			}
		}()

		wg.Wait()
		if wgErr != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		role := u.Role
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
		scope.Role = role

		ctx = jwt.SetScopeToContext(ctx, scope)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
