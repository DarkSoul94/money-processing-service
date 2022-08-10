package app

import (
	"context"

	"github.com/DarkSoul94/money-processing-service/models"
)

type Repository interface {
	CreateClient(ctx context.Context, mClient models.Client) (uint64, error)
	Close() error
}
