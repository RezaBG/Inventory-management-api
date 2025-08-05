package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	products := r.Group("/products")
	{
		products.GET("/", GetProducts)
		products.POST("/", CreateProduct)
	}
}
