package client

type Client struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Balance int    `json:"balance"`
}

type CreateClientDTO struct {
	Name    string `json:"name"`
	IP      string `json:"ip"`
	Balance int    `json:"balance"`
}
