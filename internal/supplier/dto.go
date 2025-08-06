package supplier

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
