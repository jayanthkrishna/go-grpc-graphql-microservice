package order

import (
	"context"
	"time"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}
type Order struct {
	Id         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountId  string
	Products   []OrderedProduct
}

type OrderedProduct struct {
	Id          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}

type orderService struct {
	r Repository
}

func NewService(r Repository) Service {
	return &orderService{r}
}

func (s *orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {

}
func (s *orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {

}
