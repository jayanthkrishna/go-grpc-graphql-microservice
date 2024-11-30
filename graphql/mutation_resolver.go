package main

import (
	"context"
	"errors"
	"time"

	"github.com/jayanthkrishna/go-grpc-graphql-microservice/order"
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, account AccountInput) (*Account, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := r.server.accountClient.PostAccount(ctx, account.Name)

	if err != nil {
		return nil, err
	}

	return &Account{
		Id:   res.Id,
		Name: res.Name,
	}, nil

}

func (r *mutationResolver) CreateProduct(ctx context.Context, product ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := r.server.catalogClient.PostProduct(ctx, product.Name, product.Description, product.Price)

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          res.Id,
		Name:        res.Name,
		Description: res.Description,
		Price:       res.Price,
	}, nil

}

func (r *mutationResolver) CreateOrder(ctx context.Context, inp OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	products := []order.OrderedProduct{}

	for _, p := range inp.Products {
		if p.Quantity <= 0 {
			return nil, errors.New("invalid parameter")
		}
		products = append(products, order.OrderedProduct{
			Id:       p.ID,
			Quantity: uint32(p.Quantity),
		})
	}
	res, err := r.server.orderClient.PostOrder(ctx, inp.AccountID, products)

	if err != nil {
		return nil, err
	}

	return &Order{
		ID:         res.Id,
		CreatedAt:  res.CreatedAt,
		TotalPrice: res.TotalPrice,
	}, nil

}
