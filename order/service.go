package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
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
	repository Repository
}

func NewService(r Repository) Service {
	return &orderService{r}
}

func (s *orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {

	order := &Order{
		Id:        ksuid.New().String(),
		CreatedAt: time.Now(),
		AccountId: accountID,
		Products:  products,
	}

	order.TotalPrice = 0.0

	for _, p := range products {
		order.TotalPrice += p.Price * float64(p.Quantity)

	}

	err := s.repository.PutOrder(ctx, *order)

	if err != nil {
		return nil, err
	}

	return order, nil

}
func (s *orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return s.repository.GetOrdersForAccount(ctx, accountID)
}
