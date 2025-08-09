package supplier

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type Service interface {
	CreateNewSupplier(input CreateSupplierInput) (*SupplierResponse, error)
	GetAllSuppliers() ([]SupplierResponse, error)
	GetSupplierByID(id string) (*SupplierResponse, error)
	UpdateExistingSupplier(id string, input UpdateSupplierInput) (*SupplierResponse, error)
	DeleteSupplierByID(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func toSupplierResponse(supplier Supplier) SupplierResponse {
	return SupplierResponse{
		ID:            supplier.ID,
		CreatedAt:     supplier.CreatedAt,
		Name:          supplier.Name,
		ContactPerson: supplier.ContactPerson,
		Email:         supplier.Email,
		Phone:         supplier.Phone,
	}
}

func (s *service) CreateNewSupplier(input CreateSupplierInput) (*SupplierResponse, error) {
	newSupplier := Supplier{
		Name:          input.Name,
		ContactPerson: input.ContactPerson,
		Email:         input.Email,
		Phone:         input.Phone,
	}

	// the service calls the repository to save data
	savedSupplier, err := s.repo.Save(&newSupplier)
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

	response := toSupplierResponse(*savedSupplier)
	return &response, nil

}

func (s *service) GetAllSuppliers() ([]SupplierResponse, error) {
	suppliers, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []SupplierResponse
	for _, supplier := range suppliers {
		responses = append(responses, toSupplierResponse(supplier))
	}
	return responses, nil
}

func (s *service) GetSupplierByID(id string) (*SupplierResponse, error) {
	supplier, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	response := toSupplierResponse(*supplier)
	return &response, nil
}

func (s *service) UpdateExistingSupplier(id string, input UpdateSupplierInput) (*SupplierResponse, error) {
	supplier, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update the fields
	supplier.Name = input.Name
	supplier.ContactPerson = input.ContactPerson
	supplier.Email = input.Email
	supplier.Phone = input.Phone

	updatedSupplier, err := s.repo.Update(supplier)
	if err != nil {
		return nil, err
	}

	response := toSupplierResponse(*updatedSupplier)
	return &response, nil
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
