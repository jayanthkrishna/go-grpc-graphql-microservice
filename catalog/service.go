package catalog

import "context"

type Service interface {
	PostProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type Product struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type catalogService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &catalogService{r}
}

func (s *catalogService) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	product := &Product{
		Name:        name,
		Description: description,
		Price:       price,
	}
	err := s.repository.PutProduct(ctx, *product)

	return product, err

}
func (s *catalogService) GetProduct(ctx context.Context, id string) (*Product, error) {

	return s.repository.GetProductByID(ctx, id)

}
func (s *catalogService) GetProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}

	return s.repository.ListProducts(ctx, skip, take)

}
func (s *catalogService) GetProductsByIDs(ctx context.Context, ids []string) ([]Product, error) {

	return s.repository.ListProductsWithIDs(ctx, ids)

}
func (s *catalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.SearchProducts(ctx, query, skip, take)
}
