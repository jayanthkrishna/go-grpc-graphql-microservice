package main

import (
	"context"
	"time"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := r.server.orderClient.GetOrdersForAccount(ctx, obj.Id)

	if err != nil {
		return nil, err
	}

	orders := []*Order{}

	for _, o := range orderList {

		order := &Order{
			ID:         o.Id,
			CreatedAt:  o.CreatedAt,
			TotalPrice: o.TotalPrice,
		}

		products := []*Product{}

		for _, p := range o.Products {
			products = append(products, &Product{
				ID:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}

		order.Products = products
		orders = append(orders, order)
	}

	return orders, nil

}
