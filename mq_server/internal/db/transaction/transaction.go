package transaction

import (
	"context"
	"fmt"
	"mq_server/internal/db"

	"github.com/jackc/pgconn"
)

type repository struct {
	sqlclient db.SQLClient
}

func NewRepository(client db.SQLClient) StorageRepository {
	return &repository{
		sqlclient: client,
	}
}

func (r *repository) Create(ctx context.Context, transaction *Transaction) error {
	q := `
		INSERT INTO transaction (info,sum,client_id) 
		VALUES ($1,$2,$3) 
		RETURNING id
	`
	err := r.sqlclient.QueryRow(ctx, q, transaction.Info, transaction.Sum, transaction.Client.ID).Scan(&transaction.ID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			fmt.Println(newErr)
			return newErr
		}
		return err
	}

	return nil
}

func (r *repository) FindAll(ctx context.Context) (b []Transaction, err error) {
	return nil, nil
}
func (r *repository) FindOne(ctx context.Context, id string) (Transaction, error) {
	return Transaction{}, nil
}
func (r *repository) Update(ctx context.Context, transaction Transaction) error {
	return nil
}
func (r *repository) Delete(ctx context.Context, id string) error {
	return nil
}
