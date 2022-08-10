package http

import "github.com/DarkSoul94/money-processing-service/models"

type newClient struct {
	Name string `json:"name"`
}

func (h *Handler) toModelClient(client newClient) models.Client {
	return models.Client{
		Name: client.Name,
	}
}
