package product

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetProducts(c *gin.Context) {
	products, err := h.svc.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProductByID retrieves a single product by its ID.
// @Summary      Get a single product
// @Description  Retrieves the details of a single product, including its real-time inventory count.
// @Tags         Products
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  ProductResponse
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id} [get]
func (h *Handler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.svc.GetProductByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		// Handle other errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct creates a new product.
// @Summary      Create a new product
// @Description  Adds a new product to the system. The initial quantity will be 0.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product body CreateProductInput true "Product Information"
// @Success      201  {object}  ProductResponse // <-- FIXED
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdProduct, err := h.svc.CreateNewProduct(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	c.JSON(http.StatusCreated, createdProduct)
}

// UpdateProduct updates an existing product's details.
// @Summary      Update a product
// @Description  Updates a product's details (e.g., name, price). Note: Quantity cannot be updated here.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Product ID"
// @Param        product body UpdateProductInput true "Product Update Information"
// @Success      200  {object}  ProductResponse // <-- FIXED
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id} [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var input UpdateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedProduct, err := h.svc.UpdateExistingProduct(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

// DeleteProduct deletes a product.
// @Summary      Delete a product
// @Description  Deletes a product from the system by its ID.
// @Tags         Products
// @Security     BearerAuth
// @Param        id   path      int  true  "Product ID"
// @Success      204  "No Content"
// @Failure      500  {object}  map[string]interface{}
// @Router       /products/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := h.svc.DeleteProductByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
