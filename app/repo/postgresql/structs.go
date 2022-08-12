package postgresql

import (
	"context"
	"database/sql"
	"time"

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
	Id        uuid.UUID       `db:"id"`
	CreatedAt time.Time       `db:"created_at"`
	Type      uint            `db:"type"`
	From      sql.NullInt64   `db:"from_account_id"`
	To        sql.NullInt64   `db:"to_account_id"`
	Amount    decimal.Decimal `db:"amount"`
}

func (r *postgreRepo) toDbTransaction(mTransaction models.Transaction) dbTransaction {
	transaction := dbTransaction{
		Id:        mTransaction.Id,
		CreatedAt: mTransaction.CreatedAt,
		Type:      mTransaction.Type.Id,
		Amount:    mTransaction.Amount,
	}

	if mTransaction.From.Id != 0 {
		transaction.From.Int64 = int64(mTransaction.From.Id)
		transaction.From.Valid = true
	} else {
		transaction.From.Valid = false
	}

	if mTransaction.To.Id != 0 {
		transaction.To.Int64 = int64(mTransaction.To.Id)
		transaction.To.Valid = true
	} else {
		transaction.To.Valid = false
	}

	return transaction
}

func (r *postgreRepo) toModelTransaction(transaction dbTransaction) models.Transaction {
	mTransaction := models.Transaction{
		Id:        transaction.Id,
		CreatedAt: transaction.CreatedAt,
		Type:      models.TransactionType{Id: transaction.Type},
		Amount:    transaction.Amount,
	}

	if transaction.From.Valid {
		mTransaction.From = models.Account{
			Id: uint64(transaction.From.Int64),
		}
	}

	if transaction.To.Valid {
		mTransaction.To = models.Account{
			Id: uint64(transaction.To.Int64),
		}
	}

	return mTransaction
}
