package postgresql

import (
	"context"

	"github.com/DarkSoul94/money-processing-service/models"
	"github.com/shopspring/decimal"
)

type dbClient struct {
	Id   uint64 `db:"id"`
	Name string `db:"name"`
}

func (r *postgreRepo) toDbClient(client models.Client) dbClient {
	return dbClient{
		Id:   client.Id,
		Name: client.Name,
	}
}

func (r *postgreRepo) toModelClient(client dbClient) models.Client {
	return models.Client{
		Id:   client.Id,
		Name: client.Name,
	}
}

type dbCurrency struct {
	Id   uint   `db:"id"`
	Name string `db:"name"`
}

func (r *postgreRepo) toModelCurrency(currency dbCurrency) models.Currency {
	return models.Currency{
		Id:   currency.Id,
		Name: currency.Name,
	}
}

type dbAccount struct {
	Id         uint64          `db:"id"`
	ClientID   uint64          `db:"client_id"`
	CurrencyID uint            `db:"currency_id"`
	Ballance   decimal.Decimal `db:"ballance"`
}

func (r *postgreRepo) toDbAccount(account models.Account) dbAccount {
	return dbAccount{
		Id:         account.Id,
		ClientID:   account.Client.Id,
		CurrencyID: account.Currency.Id,
		Ballance:   account.Ballance,
	}
}

func (r *postgreRepo) toModelAccount(ctx context.Context, account dbAccount) (models.Account, error) {
	client, err := r.GetClientByID(ctx, account.ClientID)
	if err != nil {
		return models.Account{}, err
	}

	currency, err := r.GetCurrencyByID(ctx, account.CurrencyID)
	if err != nil {
		return models.Account{}, err
	}

	return models.Account{
		Id:       account.Id,
		Client:   client,
		Currency: currency,
		Ballance: account.Ballance,
	}, nil
}
