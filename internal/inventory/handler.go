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

// CreateTransaction creates a new inventory transaction (e.g., stock-in, stock-out).
// @Summary      Create an inventory transaction
// @Description  Creates a new stock movement record. Use positive quantity for stock-in, negative for stock-out.
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        transaction body CreateTransactionInput true "Transaction Details"
// @Success      201  {object}  TransactionResponse //
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /inventory/transactions [post]
func (h *Handler) CreateTransaction(c *gin.Context) {
	// Authenticate the user
	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, ok := currentUser.(*user.User)
	if !ok || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user context"})
		return
	}

	var input CreateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTransaction, err := h.svc.CreateTransaction(input, *user)
	if err != nil {
		// Return a 400 Bad Request for business logic errors (e.g., negative stock-in).
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTransaction)

}
