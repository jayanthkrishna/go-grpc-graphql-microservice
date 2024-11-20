package catalog

import (
	"context"

	"github.com/jayanthkrishna/go-grpc-graphql-microservice/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	c := pb.NewCatalogServiceClient(conn)

	return &Client{
		conn:    conn,
		service: c,
	}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	r, err := c.service.PostProduct(ctx, &pb.PostProductRequest{
		Name:        name,
		Description: description,
		Price:       price,
	})

	if err != nil {
		return nil, err
	}

	return &Product{
		Id:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil

}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {

	r, err := c.service.GetProduct(ctx, &pb.GetProductRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return &Product{
		Id:          id,
		Name:        r.Product.Id,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil

}

func (c *Client) GetProducts(ctx context.Context, skip uint64, take uint64, ids []string, query string) ([]Product, error) {

	r, err := c.service.GetProducts(ctx, &pb.GetProductsRequest{
		Skip:  skip,
		Take:  take,
		Ids:   ids,
		Query: query,
	})

	if err != nil {
		return nil, err
	}

	res := []Product{}

	for _, product := range r.Products {

		res = append(res, Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return res, nil

}
