package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/user"
	"github.com/pt010104/api-golang/pkg/jwt"
	"github.com/pt010104/api-golang/pkg/log"
	"net/http"
	"os"
)

type Middleware struct {
	l    log.Logger
	repo user.Repo
}

func New(l log.Logger, repo user.Repo) Middleware {
	return Middleware{
		l:    l,
		repo: repo,
	}
}
func (mw Middleware) ResetPasswordMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Query("token")
		if token == "" {

			c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		payload, err := jwt.Verify(token, os.Getenv("SUPER_SECRET_KEY"))
		if err != nil {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if payload.Type != "reset-request" {

			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
			c.Abort()
			return
		}

		c.Set("userID", payload.UserID)

		c.Next()
	}
}
