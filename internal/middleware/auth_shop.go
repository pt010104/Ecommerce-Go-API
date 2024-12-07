package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/models"
	"github.com/pt010104/api-golang/pkg/jwt"
	"github.com/pt010104/api-golang/pkg/response"
)

func (m Middleware) AuthShop() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("x-client-id")

		ctx := c.Request.Context()

		shop, err := m.shopUC.Detail(ctx, models.Scope{
			UserID: userID,
		}, "")
		if err != nil {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		scope, ok := jwt.GetScopeFromContext(ctx)
		if !ok {
			response.Unauthorized(c)
			c.Abort()
			return
		}

		scope.ShopID = shop.ID.Hex()

		ctx = jwt.SetScopeToContext(ctx, scope)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
