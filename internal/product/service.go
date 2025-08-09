package product

type InventoryStockCalculator interface {
	CalculateStockForProduct(productID uint) (int, error)
}

type Service interface {
	GetAllProducts() ([]ProductResponse, error)
	GetProductByID(id string) (*ProductResponse, error)
	CreateNewProduct(input CreateProductInput) (*ProductResponse, error)
	UpdateExistingProduct(id string, input UpdateProductInput) (*ProductResponse, error)
	DeleteProductByID(id string) error
}

type service struct {
	productRepo     Repository
	stockCalculator InventoryStockCalculator
}

func NewService(productRepo Repository, stockCalculator InventoryStockCalculator) Service {
	return &service{
		productRepo:     productRepo,
		stockCalculator: stockCalculator,
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

func (s *service) CreateNewProduct(input CreateProductInput) (*ProductResponse, error) {
	newProduct := Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Quantity:    0,
	}

	savedProduct, err := s.productRepo.Save(&newProduct)
	if err != nil {
		return nil, err
	}

	response := &ProductResponse{
		ID:                 savedProduct.ID,
		Name:               savedProduct.Name,
		Description:        savedProduct.Description,
		Price:              savedProduct.Price,
		CalculatedQuantity: 0, // Initial quantity is always 0
	}

	return response, nil
}

func (s *service) UpdateExistingProduct(id string, input UpdateProductInput) (*ProductResponse, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if they are provided in the input.
	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price

	updatedProduct, err := s.productRepo.Update(product)
	if err != nil {
		return nil, err
	}

	// after updating, we need to recalculate the stock
	quantity, err := s.stockCalculator.CalculateStockForProduct(updatedProduct.ID)
	if err != nil {
		return nil, err
	}

	response := &ProductResponse{
		ID:                 updatedProduct.ID,
		Name:               updatedProduct.Name,
		Description:        updatedProduct.Description,
		Price:              updatedProduct.Price,
		CalculatedQuantity: quantity,
	}

	return response, nil
}

func (s *service) DeleteProductByID(id string) error {
	return s.productRepo.Delete(id)
}
