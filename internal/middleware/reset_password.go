package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/pkg/jwt"
	"github.com/pt010104/api-golang/pkg/response"
)

func (mw Middleware) ResetPasswordMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		payload, err := jwt.Verify(token, os.Getenv("SUPER_SECRET_KEY"))
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		if payload.Type != "reset-request" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("userID", payload.UserID)

		c.Next()
	}
}
func (mw Middleware) VerifyMidleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		payload, err := jwt.Verify(token, os.Getenv("SUPER_SECRET_KEY"))
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		if payload.Type != "verify" {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set("userID", payload.UserID)

		c.Next()
	}
}
