package supplier

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateSupplier(c *gin.Context) {
	var input CreateSupplierInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	supplier, err := h.svc.CreateNewSupplier(input)
	if err != nil {
		// Check if the error message is our new specific message from the service
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		// For all other errors, return a generic 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create supplier"})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

func (h *Handler) GetAllSuppliers(c *gin.Context) {
	suppliers, err := h.svc.GetAllSuppliers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch suppliers"})
		return
	}
	c.JSON(http.StatusOK, suppliers)
}

func (h *Handler) GetSupplierByID(c *gin.Context) {
	id := c.Param("id")
	supplier, err := h.svc.GetSupplierByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch supplier"})
		return
	}
	c.JSON(http.StatusOK, supplier)
}

func (h *Handler) UpdateSupplier(c *gin.Context) {
	id := c.Param("id")
	var input UpdateSupplierInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	supplier, err := h.svc.UpdateExistingSupplier(id, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier"})
		return
	}
	c.JSON(http.StatusOK, supplier)
}

func (h *Handler) DeleteSupplier(c *gin.Context) {
	id := c.Param("id")
	err := h.svc.DeleteSupplierByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete supplier"})
		return
	}
	c.Status(http.StatusNoContent)
}
