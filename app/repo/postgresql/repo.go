package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/DarkSoul94/money-processing-service/app"
	"github.com/DarkSoul94/money-processing-service/models"
	"github.com/DarkSoul94/money-processing-service/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type postgreRepo struct {
	db *sqlx.DB
}

func NewPostgreRepo(db *sql.DB) app.Repository {
	return &postgreRepo{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (r *postgreRepo) CreateClient(ctx context.Context, mClient models.Client) (uint64, error) {
	var (
		client dbClient
		query  string
		id     int64
		err    error
	)

	client = r.toDbClient(mClient)

	query = `INSERT INTO clients (name) VALUES ($1) RETURNING id`

	err = r.db.GetContext(ctx, &id, query, client.Name)
	if err != nil {
		logger.LogError(
			"Create client",
			"app/repo/postgresql/repo",
			fmt.Sprintf("name: %s", client.Name),
			err,
		)
		return 0, errors.New("failed insert client to db")
	}

	return uint64(id), nil
}

func (r *postgreRepo) GetClientByID(ctx context.Context, id uint64) (models.Client, error) {
	var (
		client dbClient
		query  string
		err    error
	)

	query = `SELECT * FROM clients WHERE id=$1`
	err = r.db.GetContext(ctx, &client, query, id)
	if err != nil {
		logger.LogError(
			"select client by id",
			"app/repo/postgresql/repo",
			fmt.Sprintf("client_id: %d", id),
			err,
		)
		return models.Client{}, errors.New("failed select client from db")
	}

	return r.toModelClient(client), nil
}

func (r *postgreRepo) GetClientAccountsID(ctx context.Context, id uint64) ([]uint64, error) {
	var (
		idList []uint64
		query  string
		err    error
	)

	query = "SELECT id FROM accounts WHERE client_id=$1"

	err = r.db.SelectContext(ctx, &idList, query, id)
	if err != nil {
		logger.LogError(
			"select client accounts id",
			"app/repo/postgresql/repo",
			fmt.Sprintf("client_id: %d", id),
			err,
		)
		return nil, errors.New("failed select client accounts id from db")
	}

	return idList, nil
}

func (r *postgreRepo) GetCurrencyByID(ctx context.Context, id uint) (models.Currency, error) {
	var (
		currency dbCurrency
		query    string
		err      error
	)

	query = `SELECT * FROM currencys WHERE id=$1`

	err = r.db.GetContext(ctx, &currency, query, id)
	if err != nil {
		logger.LogError(
			"select currency",
			"app/repo/postgresql/repo",
			fmt.Sprintf("currency_id: %d", id),
			err,
		)
		return models.Currency{}, errors.New("failed select currency from db")
	}

	return r.toModelCurrency(currency), nil
}

func (r *postgreRepo) CreateAccount(ctx context.Context, mAccount models.Account) (uint64, error) {
	var (
		account dbAccount
		query   string
		id      int64
		err     error
	)

	account = r.toDbAccount(mAccount)

	query = `INSERT INTO accounts (client_id, currency_id, ballance) VALUES ($1, $2, $3) RETURNING id`

	err = r.db.GetContext(ctx, &id, query, account.ClientID, account.CurrencyID, account.Ballance)
	if err != nil {
		logger.LogError(
			"Create account",
			"app/repo/postgresql/repo",
			fmt.Sprintf("client_id: %d, currency_id: %d, ballance: %d", account.ClientID, account.CurrencyID, account.Ballance),
			err,
		)
		return 0, errors.New("failed insert account to db")
	}

	return uint64(id), nil
}

func (r *postgreRepo) GetAccountByID(ctx context.Context, id uint64) (models.Account, error) {
	var (
		account dbAccount
		query   string
		err     error
	)

	query = `SELECT * FROM accounts WHERE id = $1`

	err = r.db.GetContext(ctx, &account, query, id)
	if err != nil {
		logger.LogError(
			"Get account",
			"app/repo/postgresql/repo",
			fmt.Sprintf("account_id: %d", id),
			err,
		)
		return models.Account{}, errors.New("failed select account from db")
	}

	return r.toModelAccount(ctx, account)
}

func (r *postgreRepo) Close() error {
	return r.db.Close()
}
