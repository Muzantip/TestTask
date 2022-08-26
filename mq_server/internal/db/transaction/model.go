package transaction

import "mq_server/internal/db/client"

type Transaction struct {
	ID     string        `json:"id"`
	Info   string        `json:"info"`
	Sum    int           `json:"sum"`
	Client client.Client `json:"client"`
}

type CreateTransactionDTO struct {
}
