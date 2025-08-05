package product

type Service interface {
	GetAllProducts() ([]Product, error)
	CreateNewProduct(input CreateProductInput) (*Product, error)
	UpdateExistingProduct(id string, input UpdateProductInput) (*Product, error)
	DeleteProductByID(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllProducts() ([]Product, error) {
	return s.repo.FindAll()
}

func (s *service) CreateNewProduct(input CreateProductInput) (*Product, error) {
	newProduct := Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Quantity:    input.Quantity,
	}
	return s.repo.Save(&newProduct)
}

func (s *service) UpdateExistingProduct(id string, input UpdateProductInput) (*Product, error) {
	// First, get the existing product.
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err // Return error if not found
	}

	// Update fields if they are provided in the input.
	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Quantity = input.Quantity

	// Save the updated product.
	return s.repo.Update(product)
}

func (s *service) DeleteProductByID(id string) error {
	return s.repo.Delete(id)
}
