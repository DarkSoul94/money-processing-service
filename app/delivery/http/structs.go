package http

import (
	"github.com/DarkSoul94/money-processing-service/models"
	"github.com/shopspring/decimal"
)

type newClient struct {
	Name string `json:"name"`
}

func (h *Handler) toModelClient(client newClient) models.Client {
	return models.Client{
		Name: client.Name,
	}
}

type outClient struct {
	Id         uint64   `json:"id"`
	Name       string   `json:"name"`
	AccountsID []uint64 `json:"accounts_id"`
}

func (h *Handler) toOutClient(client models.Client, idList []uint64) outClient {
	return outClient{
		Id:         client.Id,
		Name:       client.Name,
		AccountsID: idList,
	}
}

type newAccount struct {
	ClientID   uint64 `json:"client_id"`
	CurrencyID uint   `json:"currency_id"`
}

func (h *Handler) toModelAccount(account newAccount) models.Account {
	return models.Account{
		Client: models.Client{
			Id: account.ClientID,
		},
		Currency: models.Currency{
			Id: account.CurrencyID,
		},
		Ballance: decimal.Decimal{},
	}
}
