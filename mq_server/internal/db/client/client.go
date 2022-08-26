package client

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

func (r *repository) Create(ctx context.Context, client *Client) error {
	q := `
		INSERT INTO client (name,ip,balance) 
		VALUES ($1,$2,$3) 
		RETURNING id
	`
	err := r.sqlclient.QueryRow(ctx, q, client.Name, client.IP, client.Balance).Scan(&client.ID)
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

func (r *repository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *repository) FindAll(ctx context.Context) (b []Client, err error) {
	q := `SELECT id, name FROM client`
	rows, err := r.sqlclient.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	clients := make([]Client, 0)

	for rows.Next() {
		var cl Client
		err := rows.Scan(&cl.ID, &cl.Name, &cl.IP, &cl.Balance)
		if err != nil {
			return nil, err
		}

		clients = append(clients, cl)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *repository) FindOne(ctx context.Context, name string) (Client, error) {

	q := `SELECT id, name, ip, balance FROM client WHERE name = $1`

	var cl Client
	err := r.sqlclient.QueryRow(ctx, q, name).Scan(&cl.ID, &cl.Name, &cl.IP, &cl.Balance)
	if err != nil {
		return Client{}, err
	}
	return cl, nil
}

func (r *repository) Update(ctx context.Context, client Client) error {

	q := `UPDATE public.client SET balance = $1 WHERE name = $2`

	var cl Client
	err := r.sqlclient.QueryRow(ctx, q, client.Balance, client.Name).Scan(&cl.ID, &cl.Name, &cl.IP, &cl.Balance)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL error: %s, Detail: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			fmt.Println(newErr)
			return newErr
		}
	}

	return nil
}
