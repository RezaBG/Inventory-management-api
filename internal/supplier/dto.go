package supplier

import "time"

type CreateSupplierInput struct {
	Name          string `json:"name" binding:"required"`
	ContactPerson string `json:"contactPerson"`
	Email         string `json:"email" binding:"required,email"`
	Phone         string `json:"phone"`
}

type UpdateSupplierInput struct {
	Name          string `json:"name" binding:"required"`
	ContactPerson string `json:"contactPerson"`
	Email         string `json:"email" binding:"email"`
	Phone         string `json:"phone"`
}

type SupplierResponse struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	Name          string    `json:"name"`
	ContactPerson string    `json:"contactPerson"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
}
