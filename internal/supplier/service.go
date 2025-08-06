package supplier

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
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

	// the service calls the repository to save data
	_, err := s.repo.Save(&newSupplier)
	if err != nil {
		var pgErr *pgconn.PgError
		// errors.As checks if the error from GORM can be converted to a PgError
		if errors.As(err, &pgErr) {
			// Check if the error code is for a "unique_violation"
			if pgErr.Code == "23505" {
				// If it is, we return a much more user-friendly error
				return nil, fmt.Errorf("supplier with email '%s' already exists", newSupplier.Email)
			}
		}
		// For any other kind of database error, we return a generic message
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &newSupplier, nil

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
