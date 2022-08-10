package app

import (
	"context"

	"github.com/DarkSoul94/money-processing-service/models"
)

type Repository interface {
	CreateClient(ctx context.Context, mClient models.Client) (uint64, error)
	GetClientByID(ctx context.Context, id uint64) (models.Client, error)

	CreateAccount(ctx context.Context, mAccount models.Account) (uint64, error)
	GetClientAccountsID(ctx context.Context, id uint64) ([]uint64, error)

	Close() error
}
