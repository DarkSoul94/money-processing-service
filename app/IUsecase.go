package app

import (
	"context"

	"github.com/DarkSoul94/money-processing-service/models"
)

type Usecase interface {
	CreateClient(ctx context.Context, client models.Client) (uint64, error)
	CreateAccount(ctx context.Context, account models.Account) (uint64, error)
}
