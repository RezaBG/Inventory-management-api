package inventory

import (
	"database/sql"

	"gorm.io/gorm"
)

type Repository interface {
	Create(tx *InventoryTransaction) error
	GetTransactionsForProduct(productID uint) ([]InventoryTransaction, error)
	CalculateStockForProduct(productID uint) (int, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(tx *InventoryTransaction) error {
	return r.db.Create(tx).Error
}

func (r *repository) GetTransactionsForProduct(productID uint) ([]InventoryTransaction, error) {
	var transactions []InventoryTransaction
	err := r.db.Where("product_id = ?", productID).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *repository) CalculateStockForProduct(productID uint) (int, error) {
	var total sql.NullInt64
	err := r.db.Model(&InventoryTransaction{}).
		Where("product_id = ?", productID).
		Select("sum(quantity_change)").
		Row().
		Scan(&total)

	if err != nil {
		return 0, err
	}

	return int(total.Int64), err
}
