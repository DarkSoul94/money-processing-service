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
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) toOutClient(client models.Client) outClient {
	return outClient{
		Id:   client.Id,
		Name: client.Name,
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

type outCurrency struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) toOutCurrency(currency models.Currency) outCurrency {
	return outCurrency{
		Id:   currency.Id,
		Name: currency.Name,
	}
}

type outAccount struct {
	Id       uint64          `json:"id"`
	Client   outClient       `json:"client"`
	Currency outCurrency     `json:"currency"`
	Ballance decimal.Decimal `json:"ballance"`
}

func (h *Handler) toOutAccount(account models.Account) outAccount {
	return outAccount{
		Id:       account.Id,
		Client:   h.toOutClient(account.Client),
		Currency: h.toOutCurrency(account.Currency),
		Ballance: account.Ballance,
	}
}

type newTransaction struct {
	Type          int             `json:"type"`
	FromAccountID uint64          `json:"from_account_id,omitempty"`
	ToAccountID   uint64          `json:"to_account_id,omitempty"`
	Amount        decimal.Decimal `json:"amount"`
}

func (h *Handler) toModelTransction(transaction newTransaction) models.Transaction {
	return models.Transaction{
		Type: models.TransactionType(transaction.Type),
		From: models.Account{
			Id: transaction.FromAccountID,
		},
		To: models.Account{
			Id: transaction.ToAccountID,
		},
		Amount: transaction.Amount,
	}
}
