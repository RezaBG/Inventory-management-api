package inventory

import (
	"net/http"

	"github.com/RezaBG/Inventory-management-api/internal/user"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	// Authenticate the user
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := currentUser.(user.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
		return
	}

	var input CreateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTransaction, err := h.svc.CreateTransaction(input, user)
	if err != nil {
		// Return a 400 Bad Request for business logic errors (e.g., negative stock-in).
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTransaction)

}
