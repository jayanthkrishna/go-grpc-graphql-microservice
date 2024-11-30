package order

import (
	"context"
	"log"
	"time"

	"github.com/jayanthkrishna/go-grpc-graphql-microservice/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	c := pb.NewOrderServiceClient(conn)

	return &Client{
		conn:    conn,
		service: c,
	}, nil
}

func (c *Client) PostOrder(ctx context.Context, accountId string, products []OrderedProduct) (*Order, error) {
	protoProducts := []*pb.PostOrderRequest_OrderProduct{}

	for _, p := range products {
		protoProducts = append(protoProducts, &pb.PostOrderRequest_OrderProduct{
			ProductId: p.Id,
			Quantity:  p.Quantity,
		})
	}

	res, err := c.service.PostOrder(ctx, &pb.PostOrderRequest{
		AccountId: accountId,
		Products:  protoProducts,
	})

	if err != nil {
		return nil, err
	}

	newOrder := res.Order
	newOrderCreatedAt := time.Time{}

	newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)
	return &Order{
		Id:         newOrder.Id,
		AccountId:  newOrder.AccountId,
		TotalPrice: newOrder.TotalPrice,
		CreatedAt:  newOrderCreatedAt,
	}, nil
}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	r, err := c.service.GetOrdersForAccount(ctx, &pb.GetOrdersForAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	orders := []Order{}
	for _, orderProto := range r.Orders {
		newOrder := Order{
			Id:         orderProto.Id,
			TotalPrice: orderProto.TotalPrice,
			AccountId:  orderProto.AccountId,
		}
		newOrder.CreatedAt = time.Time{}
		newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)

		products := []OrderedProduct{}
		for _, p := range orderProto.Products {
			products = append(products, OrderedProduct{
				Id:          p.Id,
				Quantity:    p.Quantity,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
		newOrder.Products = products

		orders = append(orders, newOrder)
	}
	return orders, nil
}
