package app

import (
	"context"

	"github.com/DarkSoul94/money-processing-service/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Repository interface {
	CreateClient(ctx context.Context, mClient models.Client) (uint64, error)
	GetClientByID(ctx context.Context, id uint64) (models.Client, error)

	CreateAccount(ctx context.Context, mAccount models.Account, clientID uint64) (uint64, error)
	GetClientAccountsID(ctx context.Context, id uint64) ([]uint64, error)
	GetAccountByID(ctx context.Context, id uint64) (models.Account, error)

	UpdateBalance(ctx context.Context, transactionType models.TransactionType, accountID uint64, amount decimal.Decimal) error
	TransferMoney(ctx context.Context, fromAccountID uint64, toAccountID uint64, amount decimal.Decimal) error

	CreateTransaction(ctx context.Context, mTransaction models.Transaction) (uuid.UUID, error)
	GetTransactionsListByAccountID(ctx context.Context, accountID uint64) ([]models.Transaction, error)

	Close() error
}
