package order

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Repository interface {
	Close()

	PutOrder(ctx context.Context, o Order) error

	GetOrdersForAccount(ctx context.Context, id string) ([]Order, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return &postgresRepository{db}, nil

}

func (r *postgresRepository) Close() {
	r.db.Close()

}

func (r *postgresRepository) PutOrder(ctx context.Context, o Order) error {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()
	_, err = tx.ExecContext(ctx,
		"INSERT INTO ORDERS(id,created_at,account_id,total_price) VALUES ($1,$2,$3,$4)",
		o.Id, o.CreatedAt, o.AccountId, o.TotalPrice,
	)

	if err != nil {
		return err
	}

	stmt, _ := tx.PrepareContext(ctx, pq.CopyIn("order_products", "order_id", "product_id", "quantity"))

	for _, p := range o.Products {
		_, err = stmt.ExecContext(ctx, o.Id, p.Id, p.quantity)

		if err != nil {
			return err
		}

	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}
	stmt.Close()

	return nil

}

func (r *postgresRepository) GetOrdersForAccount(ctx context.Context, id string) ([]Order, error) {

}
