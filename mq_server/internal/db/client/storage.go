package client

import "context"

type StorageRepository interface {
	Create(ctx context.Context, client *Client) error
	FindAll(ctx context.Context) (b []Client, err error)
	FindOne(ctx context.Context, id string) (Client, error)
	Update(ctx context.Context, client Client) error
	Delete(ctx context.Context, id string) error
}
