package order

import "time"

type Service interface {
	PostOrder()
	GetOrdersForAccount()
}
type Order struct {
	Id         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountId  string
	Products   []OrderedProduct
}

type OrderedProduct struct {
}
