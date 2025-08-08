package product

type CreateProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
}

type UpdateProductInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"gt=0"`
}

type ProductResponse struct {
	ID                 uint    `json:"id"`
	Name               string  `json:"name"`
	Description        string  `json:"description"`
	Price              float64 `json:"price"`
	CalculatedQuantity int     `json:"quantity"`
}
