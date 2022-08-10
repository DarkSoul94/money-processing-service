package postgresql

import (
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
