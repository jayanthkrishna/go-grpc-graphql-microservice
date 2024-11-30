package order

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jayanthkrishna/go-grpc-graphql-microservice/account"
	"github.com/jayanthkrishna/go-grpc-graphql-microservice/catalog"
	"github.com/jayanthkrishna/go-grpc-graphql-microservice/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
	pb.UnimplementedOrderServiceServer
}

func ListenGrpc(s Service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	})
	reflection.Register(serv)

	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {

	_, err := s.accountClient.GetAccount(ctx, r.AccountId)

	if err != nil {
		log.Println("Error getting account : ", err)
		return nil, err
	}

	productIds := []string{}

	for _, p := range r.Products {
		productIds = append(productIds, p.ProductId)
	}

	orderedProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIds, "")

	if err != nil {
		log.Println("Error getting products : ,", err)
		return nil, err
	}

	products := []OrderedProduct{}

	for _, p := range orderedProducts {
		product := OrderedProduct{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    0,
		}

		for _, rp := range r.Products {
			if rp.ProductId == product.Id {
				product.Quantity = rp.Quantity
				break
			}

			if product.Quantity != 0 {
				products = append(products, product)
			}
		}
	}

	order, err := s.service.PostOrder(ctx, r.AccountId, products)

	if err != nil {
		log.Println("Error posting order : ", err)

		return nil, err
	}

	orderProto := &pb.Order{
		Id:        order.Id,
		AccountId: order.AccountId,

		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})

	}

	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil

}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, r *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {

	accountorders, err := s.service.GetOrdersForAccount(ctx, r.AccountId)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	productIds := []string{}
	for _, order := range accountorders {
		for _, p := range order.Products {
			productIds = append(productIds, p.Id)

		}

	}

	// products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIds, "")

	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	orders := []*pb.Order{}

	for _, o := range accountorders {
		op := &pb.Order{
			Id:         o.Id,
			AccountId:  o.AccountId,
			TotalPrice: o.TotalPrice,
			Products:   []*pb.Order_OrderProduct{},
		}

		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		for _, p := range o.Products {
			op.Products = append(op.Products, &pb.Order_OrderProduct{
				Id:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    p.Quantity,
			})
		}

		orders = append(orders, op)
	}

	return &pb.GetOrdersForAccountResponse{
		Orders: orders,
	}, nil

}
