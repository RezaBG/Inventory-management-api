package inventory

import (
	"fmt"

	"github.com/RezaBG/Inventory-management-api/internal/product"
	"github.com/RezaBG/Inventory-management-api/internal/user"
)

type Service interface {
	CreateTransaction(input CreateTransactionInput, currentUser user.User) (*TransactionResponse, error)
}

type service struct {
	inventoryRepo Repository
	productRepo   product.Repository
}

func NewService(inventoryRepo Repository, productRepo product.Repository) Service {
	return &service{
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
	}
}

func (s *service) CreateTransaction(input CreateTransactionInput, currentUser user.User) (*TransactionResponse, error) {
	switch input.Type {
	case StockIn:
		if input.QuantityChange <= 0 {
			return nil, fmt.Errorf("stock-in quantity must be positive")
		}
	case StockOut:
		if input.QuantityChange >= 0 {
			return nil, fmt.Errorf("stock-out quantity must be negative")
		}
	case Adjustment:
		if input.QuantityChange == 0 {
			return nil, fmt.Errorf("quantity change for adjustment cannot be zero")
		}
	default:
		return nil, fmt.Errorf("invalid transaction type")
	}

	// Check if the product exists.
	_, err := s.productRepo.FindByID(fmt.Sprint(input.ProductID))
	if err != nil {
		return nil, fmt.Errorf("product with ID %d not found", input.ProductID)
	}

	newTransaction := &InventoryTransaction{
		ProductID:      input.ProductID,
		UserID:         currentUser.ID,
		Type:           input.Type,
		QuantityChange: input.QuantityChange,
		Notes:          input.Notes,
	}

	err = s.inventoryRepo.Create(newTransaction)
	if err != nil {
		return nil, fmt.Errorf("could not create transaction: %w", err)
	}

	// TODO: Business Rule 3: Check if a stock-out would result in negative inventory.

	// All checks passed, create the transaction.
	response := &TransactionResponse{
		ID:             newTransaction.ID,
		CreatedAt:      newTransaction.CreatedAt,
		ProductID:      newTransaction.ProductID,
		UserID:         newTransaction.UserID,
		Type:           newTransaction.Type,
		QuantityChange: newTransaction.QuantityChange,
		Notes:          newTransaction.Notes,
	}

	// Call the repository to save the transaction.
	err = s.inventoryRepo.Create(newTransaction)
	if err != nil {
		return nil, fmt.Errorf("could not save transaction: %w", err)
	}

	return response, nil
}
