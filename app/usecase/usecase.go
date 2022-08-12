package usecase

import (
	"context"
	"errors"

	"github.com/DarkSoul94/money-processing-service/app"
	"github.com/DarkSoul94/money-processing-service/models"
	"github.com/google/uuid"
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

func (u *usecase) GetClientByID(ctx context.Context, id uint64) (models.Client, []uint64, error) {
	client, err := u.repo.GetClientByID(ctx, id)
	if err != nil {
		return models.Client{}, nil, err
	}

	idList, err := u.repo.GetClientAccountsID(ctx, id)
	if err != nil {
		return models.Client{}, nil, err
	}

	return client, idList, nil
}

func (u *usecase) CreateAccount(ctx context.Context, account models.Account) (uint64, error) {
	return u.repo.CreateAccount(ctx, account)
}

func (u *usecase) GetAccountByID(ctx context.Context, id uint64) (models.Account, error) {
	return u.repo.GetAccountByID(ctx, id)
}

func (u *usecase) CreateTransaction(ctx context.Context, transaction models.Transaction) (uuid.UUID, error) {
	switch transaction.Type {
	case models.Deposit:
		return u.depositMoney(ctx, transaction)
	case models.Withdraw:
		return u.withdrawMoney(ctx, transaction)
	case models.Transfer:
		return u.transferMoney(ctx, transaction)
	}

	return uuid.UUID{}, errors.New("invalid transaction type")
}

func (u *usecase) depositMoney(ctx context.Context, transaction models.Transaction) (uuid.UUID, error) {
	transaction.From = transaction.To

	err := u.repo.UpdateBalance(ctx, transaction.Type, transaction.To.Id, transaction.Amount)
	if err != nil {
		return uuid.UUID{}, err
	}

	return u.repo.CreateTransaction(ctx, transaction)
}

func (u *usecase) withdrawMoney(ctx context.Context, transaction models.Transaction) (uuid.UUID, error) {
	transaction.To = transaction.From

	account, err := u.repo.GetAccountByID(ctx, transaction.From.Id)
	if err != nil {
		return uuid.UUID{}, err
	}

	if res := account.Ballance.Cmp(transaction.Amount); res == -1 {
		return uuid.UUID{}, errors.New("not enough money")
	}

	err = u.repo.UpdateBalance(ctx, transaction.Type, transaction.From.Id, transaction.Amount)
	if err != nil {
		return uuid.UUID{}, err
	}

	return u.repo.CreateTransaction(ctx, transaction)
}

func (u *usecase) transferMoney(ctx context.Context, transaction models.Transaction) (uuid.UUID, error) {
	fromAccount, err := u.repo.GetAccountByID(ctx, transaction.From.Id)
	if err != nil {
		return uuid.UUID{}, err
	}

	toAccount, err := u.repo.GetAccountByID(ctx, transaction.To.Id)
	if err != nil {
		return uuid.UUID{}, err
	}

	if fromAccount.Currency.Id != toAccount.Currency.Id {
		return uuid.UUID{}, errors.New("accounts with different currencies")
	}

	if res := fromAccount.Ballance.Cmp(transaction.Amount); res == -1 {
		return uuid.UUID{}, errors.New("not enough money")
	}

	err = u.repo.TransferMoney(ctx, transaction.From.Id, transaction.To.Id, transaction.Amount)
	if err != nil {
		return uuid.UUID{}, err
	}

	return u.repo.CreateTransaction(ctx, transaction)
}

func (u *usecase) GetTransactionsListByAccountID(ctx context.Context, accountID uint64) ([]models.Transaction, error) {
	return u.repo.GetTransactionsListByAccountID(ctx, accountID)
}
