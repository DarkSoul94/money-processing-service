package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/DarkSoul94/money-processing-service/app"
	"github.com/DarkSoul94/money-processing-service/models"
	"github.com/DarkSoul94/money-processing-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
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

func (r *postgreRepo) CreateAccount(ctx context.Context, mAccount models.Account, clientID uint64) (uint64, error) {
	var (
		account dbAccount
		query   string
		id      int64
		err     error
	)

	account = r.toDbAccount(mAccount, clientID)

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

func (r *postgreRepo) UpdateBalance(ctx context.Context, accountID uint64, newBalance decimal.Decimal) error {
	var (
		query string
		err   error
	)

	query = `UPDATE accounts SET ballance = $1 WHERE id = $2`

	_, err = r.db.ExecContext(ctx, query, newBalance, accountID)
	if err != nil {
		logger.LogError(
			"Update balance",
			"app/repo/postgresql/repo",
			fmt.Sprintf("account_id: %d, new balance: %d", accountID, newBalance),
			err,
		)
		return errors.New("failed update ballance in db")
	}

	return nil
}

func (r *postgreRepo) TransferMoney(ctx context.Context, fromAccountID uint64, toAccountID uint64, amount decimal.Decimal) error {
	var (
		query string
		err   error
	)

	tx, err := r.db.Begin()
	if err != nil {
		logger.LogError(
			"Begin transaction",
			"app/repo/postgresql/repo",
			"",
			err,
		)
		return errors.New("failed begin transaction")
	}

	query = `UPDATE accounts SET ballance = ballance - $1 WHERE id = $2;`
	_, err = tx.ExecContext(ctx, query, amount, fromAccountID)
	if err != nil {
		logger.LogError(
			"Withdraw money",
			"app/repo/postgresql/repo",
			fmt.Sprintf("from account id: %d, amount: %d", fromAccountID, amount),
			err,
		)

		err = tx.Rollback()
		if err != nil {
			logger.LogError(
				"Rollback transaction",
				"app/repo/postgresql/repo",
				"",
				err,
			)
			return errors.New("failed rollback transaction")
		}

		return errors.New("failed withdraw money")
	}

	query = `UPDATE accounts SET ballance = ballance + $1 WHERE id = $2;`
	_, err = tx.ExecContext(ctx, query, amount, toAccountID)
	if err != nil {
		logger.LogError(
			"Deposit money",
			"app/repo/postgresql/repo",
			fmt.Sprintf("to account id: %d, amount: %d", toAccountID, amount),
			err,
		)

		err = tx.Rollback()
		if err != nil {
			logger.LogError(
				"Rollback transaction",
				"app/repo/postgresql/repo",
				"",
				err,
			)
			return errors.New("failed rollback transaction")
		}

		return errors.New("failed deposit money")
	}

	err = tx.Commit()
	if err != nil {
		logger.LogError(
			"Commit transaction",
			"app/repo/postgresql/repo",
			"",
			err,
		)
		return errors.New("failed commit transaction")
	}

	return nil
}

func (r *postgreRepo) CreateTransaction(ctx context.Context, mTransaction models.Transaction) (uuid.UUID, error) {
	var (
		id    uuid.UUID
		query string
		err   error
	)

	query = `INSERT INTO transactions (created_at, type, from_account_id, to_account_id, amount) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	dbTransaction := r.toDbTransaction(mTransaction)

	err = r.db.GetContext(ctx, &id, query, dbTransaction.CreatedAt, dbTransaction.Type, dbTransaction.From, dbTransaction.To, dbTransaction.Amount)
	if err != nil {
		logger.LogError(
			"Create transaction",
			"app/repo/postgresql/repo",
			fmt.Sprintf("type: %d,from: %d,to: %d,value: %d,", dbTransaction.Type, dbTransaction.From.Int64, dbTransaction.To.Int64, dbTransaction.Amount),
			err,
		)
		return uuid.Nil, errors.New("failed insert transaction to db")
	}

	return id, nil
}

func (r *postgreRepo) GetTransactionsListByAccountID(ctx context.Context, accountID uint64) ([]models.Transaction, error) {
	var (
		dbTransactionsList = make([]dbTransaction, 0)
		mTransactionsList  = make([]models.Transaction, 0)
		query              string
		err                error
	)

	query = `SELECT * FROM transactions WHERE from_account_id = $1 OR to_account_id = $1`

	err = r.db.SelectContext(ctx, &dbTransactionsList, query, accountID)
	if err != nil {
		logger.LogError(
			"Get transactions list",
			"app/repo/postgresql/repo",
			fmt.Sprintf("account id: %d", accountID),
			err,
		)
		return nil, errors.New("failed select transactions list from db")
	}

	for _, dbTransaction := range dbTransactionsList {
		mTransactionsList = append(mTransactionsList, r.toModelTransaction(dbTransaction))
	}

	return mTransactionsList, nil
}

func (r *postgreRepo) Close() error {
	return r.db.Close()
}
