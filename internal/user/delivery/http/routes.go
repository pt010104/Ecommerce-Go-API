package http

import (
	"github.com/gin-gonic/gin"
	"github.com/pt010104/api-golang/internal/middleware"
)

func MapRouters(r *gin.RouterGroup, h Handler, mw middleware.Middleware) {
	r.POST("/sign-up", h.SignUp)
	r.POST("/sign-in", h.SignIn)
	r.POST("/forget-password", h.ForgetPasswordRequest)
	r.POST("/verify-request", h.VerifyEmail)
	r.POST("/distribute-new-token", h.DistributeNewToken)
	r.POST("/reset-password", mw.ResetPasswordMiddleware(), h.ResetPassword)
	r.POST("/verify", mw.VerifyMidleware(), h.VerifyUser)
	r.Use(mw.Auth())
	r.GET("/:id", h.Detail)
	r.POST("/sign-out", h.SignOut)
	r.PUT("", h.Update)

	r.POST("/address", h.AddAddress)
	r.PATCH("/address", h.UpdateAddress)
}
