package product

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Product, error)
	FindByID(id string) (*Product, error)
	Save(product *Product) (*Product, error)
	Update(product *Product) (*Product, error)
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Product, error) {
	var products []Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *repository) FindByID(id string) (*Product, error) {
	var product Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *repository) Save(product *Product) (*Product, error) {
	err := r.db.Create(product).Error
	return product, err
}

func (r *repository) Update(product *Product) (*Product, error) {
	err := r.db.Save(product).Error
	return product, err
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&Product{}, id).Error
}
