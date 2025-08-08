package product

type InventoryStockCalculator interface {
	CalculateStockForProduct(productID uint) (int, error)
}

type Service interface {
	GetAllProducts() ([]ProductResponse, error)
	GetProductByID(id string) (*ProductResponse, error)
	CreateNewProduct(input CreateProductInput) (*Product, error)
	UpdateExistingProduct(id string, input UpdateProductInput) (*Product, error)
	DeleteProductByID(id string) error
}

type service struct {
	productRepo     Repository
	stockCalculator InventoryStockCalculator
}

func NewService(productRepo Repository, inventoryRepo InventoryStockCalculator) Service {
	return &service{
		productRepo:     productRepo,
		stockCalculator: inventoryRepo,
	}
}

func (s *service) GetAllProducts() ([]ProductResponse, error) {
	products, err := s.productRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Create a slice of our response DTO
	var productResponses []ProductResponse

	// Loop through each product and calculate its stock
	for _, p := range products {
		quantity, err := s.stockCalculator.CalculateStockForProduct(p.ID)
		if err != nil {
			return nil, err
		}

		productResponses = append(productResponses, ProductResponse{
			ID:                 p.ID,
			Name:               p.Name,
			Description:        p.Description,
			Price:              p.Price,
			CalculatedQuantity: quantity,
		})
	}

	return productResponses, nil

}

func (s *service) GetProductByID(id string) (*ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	quantity, err := s.stockCalculator.CalculateStockForProduct(product.ID)
	if err != nil {
		return nil, err
	}

	response := &ProductResponse{
		ID:                 product.ID,
		Name:               product.Name,
		Description:        product.Description,
		Price:              product.Price,
		CalculatedQuantity: quantity,
	}

	return response, nil
}

func (s *service) CreateNewProduct(input CreateProductInput) (*Product, error) {
	newProduct := Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Quantity:    0,
	}

	return s.productRepo.Save(&newProduct)
}

func (s *service) UpdateExistingProduct(id string, input UpdateProductInput) (*Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if they are provided in the input.
	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price

	// Save the updated product.
	return s.productRepo.Update(product)
}

func (s *service) DeleteProductByID(id string) error {
	return s.productRepo.Delete(id)
}
