package inventory

type CreateTransactionInput struct {
	ProductID      uint            `json:"productID" binding:"required"`
	Type           TransactionType `json:"type" binding:"required,oneof=stock_in stock_out adjustment"`
	QuantityChange int             `json:"quantityChange" binding:"required"`
	Notes          string          `json:"notes,omitempty"`
}
