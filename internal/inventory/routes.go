package inventory

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h *Handler) {
	inventoryRoutes := router.Group("/inventory")
	{
		inventoryRoutes.POST("/transactions", h.CreateTransaction)
	}
}
