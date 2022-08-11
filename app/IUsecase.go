package app

import (
	"context"

	"github.com/DarkSoul94/money-processing-service/models"
	"github.com/google/uuid"
)

type Usecase interface {
	CreateClient(ctx context.Context, client models.Client) (uint64, error)
	GetClientByID(ctx context.Context, id uint64) (models.Client, []uint64, error)

	CreateAccount(ctx context.Context, account models.Account) (uint64, error)
	GetAccountByID(ctx context.Context, id uint64) (models.Account, error)

	CreateTransaction(ctx context.Context, transaction models.Transaction) (uuid.UUID, error)
}
