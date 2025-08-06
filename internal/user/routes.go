package user

import "github.com/gin-gonic/gin"

func RegisterAuthRoutes(router *gin.Engine, h *Handler) {
	router.POST("/register", h.CreateUser)
	router.POST("/login", h.Login)
	router.POST("/refresh-token", h.RefreshToken)
}
