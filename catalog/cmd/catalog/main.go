package main

import (
	"log"
	"time"

	"github.com/jayanthkrishna/go-grpc-graphql-microservice/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config

	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	var r catalog.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)

		if err != nil {
			return err
		}

		return nil

	})

	defer r.Close()

	log.Println("Listening on port 8080")

	s := catalog.NewService(r)

	log.Fatal(catalog.ListenGrpc(s, 8080))

}
