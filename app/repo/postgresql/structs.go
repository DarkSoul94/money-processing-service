package postgresql

import (
	"context"

	"github.com/DarkSoul94/money-processing-service/models"
	"github.com/google/uuid"
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

func (r *postgreRepo) toDbAccount(account models.Account, clientID uint64) dbAccount {
	return dbAccount{
		Id:         account.Id,
		ClientID:   clientID,
		CurrencyID: account.Currency.Id,
		Ballance:   account.Ballance,
	}
}

func (r *postgreRepo) toModelAccount(ctx context.Context, account dbAccount) (models.Account, error) {
	currency, err := r.GetCurrencyByID(ctx, account.CurrencyID)
	if err != nil {
		return models.Account{}, err
	}

	return models.Account{
		Id:       account.Id,
		Currency: currency,
		Ballance: account.Ballance,
	}, nil
}

type dbTransaction struct {
	Id     uuid.UUID       `db:"id"`
	Type   int             `db:"type"`
	From   uint64          `db:"from_account_id"`
	To     uint64          `db:"to_account_id"`
	Amount decimal.Decimal `db:"amount"`
}

func (r *postgreRepo) toDbTransaction(transaction models.Transaction) dbTransaction {
	return dbTransaction{
		Id:     transaction.Id,
		Type:   int(transaction.Type),
		From:   transaction.From.Id,
		To:     transaction.To.Id,
		Amount: transaction.Amount,
	}
}

func (r *postgreRepo) toModelTransaction(transaction dbTransaction) models.Transaction {
	return models.Transaction{
		Id:   transaction.Id,
		Type: models.TransactionType(transaction.Type),
		From: models.Account{
			Id: transaction.From,
		},
		To: models.Account{
			Id: transaction.To,
		},
		Amount: transaction.Amount,
	}
}
