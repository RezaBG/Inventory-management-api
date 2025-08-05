package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, h *Handler) {
	productRoutes := router.Group("/products")
	{
		productRoutes.POST("", h.CreateProduct)
		productRoutes.GET("", h.GetProducts)
		productRoutes.GET("/:id", h.GetProductByID)
		productRoutes.PUT("/:id", h.UpdateProduct)
		productRoutes.DELETE("/:id", h.DeleteProduct)
	}
}
