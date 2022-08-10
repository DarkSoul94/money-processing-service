package usecase

import (
	"context"

	"github.com/DarkSoul94/money-processing-service/app"
	"github.com/DarkSoul94/money-processing-service/models"
)

type usecase struct {
	repo app.Repository
}

func NewUsecase(repo app.Repository) app.Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) CreateClient(ctx context.Context, client models.Client) (uint64, error) {
	return u.repo.CreateClient(ctx, client)
}

func (u *usecase) CreateAccount(ctx context.Context, account models.Account) (uint64, error) {
	return u.repo.CreateAccount(ctx, account)
}
