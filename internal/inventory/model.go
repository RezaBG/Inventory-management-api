package inventory

import (
	"github.com/RezaBG/Inventory-management-api/internal/user"

	"gorm.io/gorm"
)

type TransactionType string

const (
	StockIn    TransactionType = "stock_in"
	StockOut   TransactionType = "stock_out"
	Adjustment TransactionType = "adjustment"
)

type InventoryTransaction struct {
	gorm.Model
	ProductID      uint            `json:"productID" gorm:"not null"`
	UserID         uint            `json:"userID" gorm:"not null"`
	User           user.User       `json:"user"`
	Type           TransactionType `json:"type" gorm:"not null"`
	QuantityChange int             `json:"quantityChange" gorm:"not null"`
	Notes          string          `json:"notes,omitempty"`
}
