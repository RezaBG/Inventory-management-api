package supplier

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Service interface {
	CreateNewSupplier(input CreateSupplierInput) (*Supplier, error)
	GetAllSuppliers() ([]Supplier, error)
	GetSupplierByID(id string) (*Supplier, error)
	UpdateExistingSupplier(id string, input UpdateSupplierInput) (*Supplier, error)
	DeleteSupplierByID(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateNewSupplier(input CreateSupplierInput) (*Supplier, error) {
	newSupplier := Supplier{
		Name:          input.Name,
		ContactPerson: input.ContactPerson,
		Email:         input.Email,
		Phone:         input.Phone,
	}
	return s.repo.Save(&newSupplier)
}

func (s *service) GetAllSuppliers() ([]Supplier, error) {
	return s.repo.FindAll()
}

func (s *service) GetSupplierByID(id string) (*Supplier, error) {
	return s.repo.FindByID(id)
}

func (s *service) UpdateExistingSupplier(id string, input UpdateSupplierInput) (*Supplier, error) {
	supplier, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update the fields
	supplier.Name = input.Name
	supplier.ContactPerson = input.ContactPerson
	supplier.Email = input.Email
	supplier.Phone = input.Phone

	// Call the repository's Update method
	return s.repo.Update(supplier)
}

func (s *service) DeleteSupplierByID(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("supplier with ID %s not found", id)
		}
		return err
	}

	return s.repo.Delete(id)
}
