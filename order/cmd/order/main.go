package main

import (
	"log"
	"time"

	"github.com/jayanthkrishna/go-grpc-graphql-microservice/order"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
	AccountURL  string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL  string `envconfig:"CATALOG_SERVICE_URL"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	var r order.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = order.NewPostgresRepository(cfg.DatabaseURL)

		if err != nil {
			return err
		}

		return nil

	})

	defer r.Close()

	log.Println("Listening on port 8080")

	s := order.NewService(r)

	log.Fatal(order.ListenGrpc(s, 8080))

}
