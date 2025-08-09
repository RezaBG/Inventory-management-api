package inventory

import "time"

type CreateTransactionInput struct {
	ProductID      uint            `json:"productID" binding:"required"`
	Type           TransactionType `json:"type" binding:"required,oneof=stock_in stock_out adjustment"`
	QuantityChange int             `json:"quantityChange" binding:"required"`
	Notes          string          `json:"notes,omitempty"`
}

type TransactionResponse struct {
	ID             uint            `json:"id"`
	CreatedAt      time.Time       `json:"createdAt"`
	ProductID      uint            `json:"productID"`
	UserID         uint            `json:"userID"`
	Type           TransactionType `json:"type"`
	QuantityChange int             `json:"quantityChange"`
	Notes          string          `json:"notes,omitempty"`
}
