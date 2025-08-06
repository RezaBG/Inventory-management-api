package supplier

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, h *Handler) {
	supplierRoutes := router.Group("/suppliers")
	{
		supplierRoutes.POST("/", h.CreateSupplier)
		supplierRoutes.GET("/", h.GetAllSuppliers)
		supplierRoutes.GET("/:id", h.GetSupplierByID)
		supplierRoutes.PUT("/:id", h.UpdateSupplier)
		supplierRoutes.DELETE("/:id", h.DeleteSupplier)
	}
}
