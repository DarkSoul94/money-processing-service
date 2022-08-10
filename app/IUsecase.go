package app

import (
	"context"

	"github.com/DarkSoul94/money-processing-service/models"
)

type Usecase interface {
	CreateClient(ctx context.Context, client models.Client) (uint64, error)
}
