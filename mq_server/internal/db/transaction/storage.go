package transaction

import "context"

type StorageRepository interface {
	Create(ctx context.Context, client *Transaction) error
	FindAll(ctx context.Context) (b []Transaction, err error)
	FindOne(ctx context.Context, id string) (Transaction, error)
	Update(ctx context.Context, client Transaction) error
	Delete(ctx context.Context, id string) error
}
