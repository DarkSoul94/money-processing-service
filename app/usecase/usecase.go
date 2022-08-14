package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/DarkSoul94/money-processing-service/app"
	"github.com/DarkSoul94/money-processing-service/models"
	"github.com/DarkSoul94/money-processing-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

func (u *usecase) CreateAccount(ctx context.Context, currencyID uint, clientID uint64) (uint64, error) {
	account := models.Account{
		Currency: models.Currency{Id: currencyID},
		Ballance: decimal.Decimal{},
	}
	return u.repo.CreateAccount(ctx, account, clientID)
}

func (u *usecase) GetAccountByID(ctx context.Context, id uint64) (models.Account, error) {
	return u.repo.GetAccountByID(ctx, id)
}

func (u *usecase) CreateTransaction(ctx context.Context, transaction models.Transaction) (uuid.UUID, error) {
	transaction.CreatedAt = time.Now()

	switch transaction.Type.Id {
	case models.Deposit.Id:
		return u.depositMoney(ctx, transaction)
	case models.Withdraw.Id:
		return u.withdrawMoney(ctx, transaction)
	case models.Transfer.Id:
		return u.transferMoney(ctx, transaction)
	}

	return uuid.UUID{}, errInvalidTransactionType
}

func (u *usecase) depositMoney(ctx context.Context, transaction models.Transaction) (uuid.UUID, error) {
	if transaction.To.Id == 0 {
		logger.LogError(
			"Deposit money",
			"app/usecase/usecase",
			"",
			errZeroAccountID,
		)
		return uuid.UUID{}, errZeroAccountID
	}

	transaction.From = models.Account{
		Id: 0,
	}

	account, err := u.repo.GetAccountByID(ctx, transaction.To.Id)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = u.repo.UpdateBalance(ctx, transaction.To.Id, account.Ballance.Add(transaction.Amount))
	if err != nil {
		return uuid.UUID{}, err
	}

	return u.repo.CreateTransaction(ctx, transaction)
}

func (u *usecase) withdrawMoney(ctx context.Context, transaction models.Transaction) (uuid.UUID, error) {
	if transaction.From.Id == 0 {
		logger.LogError(
			"Withdraw money",
			"app/usecase/usecase",
			"",
			errZeroAccountID,
		)
		return uuid.UUID{}, errZeroAccountID
	}

	transaction.To = models.Account{
		Id: 0,
	}

	account, err := u.repo.GetAccountByID(ctx, transaction.From.Id)
	if err != nil {
		return uuid.UUID{}, err
	}

	if res := account.Ballance.Cmp(transaction.Amount); res == -1 {
		logger.LogError(
			"Withdraw money",
			"app/usecase/usecase",
			"",
			errNotMoney,
		)
		return uuid.UUID{}, errNotMoney
	}

	err = u.repo.UpdateBalance(ctx, transaction.From.Id, account.Ballance.Sub(transaction.Amount))
	if err != nil {
		return uuid.UUID{}, err
	}

	return u.repo.CreateTransaction(ctx, transaction)
}

func (u *usecase) transferMoney(ctx context.Context, transaction models.Transaction) (uuid.UUID, error) {
	if transaction.From.Id == 0 || transaction.To.Id == 0 {
		logger.LogError(
			"Transfer money",
			"app/usecase/usecase",
			"",
			errZeroAccountID,
		)
		return uuid.UUID{}, errZeroAccountID
	}

	if transaction.From.Id == transaction.To.Id {
		logger.LogError(
			"Transfer money",
			"app/usecase/usecase",
			fmt.Sprintf("accounts id: %d", transaction.From.Id),
			errSameAccount,
		)

		return uuid.UUID{}, errSameAccount
	}

	fromAccount, err := u.repo.GetAccountByID(ctx, transaction.From.Id)
	if err != nil {
		return uuid.UUID{}, err
	}

	toAccount, err := u.repo.GetAccountByID(ctx, transaction.To.Id)
	if err != nil {
		return uuid.UUID{}, err
	}

	if fromAccount.Currency.Id != toAccount.Currency.Id {
		logger.LogError(
			"Transfer money",
			"app/usecase/usecase",
			fmt.Sprintf("from account currency id: %d, to account currency id: %d", fromAccount.Currency.Id, toAccount.Currency.Id),
			errDifferentCurrencies,
		)

		return uuid.UUID{}, errDifferentCurrencies
	}

	if res := fromAccount.Ballance.Cmp(transaction.Amount); res == -1 {
		logger.LogError(
			"Withdraw money",
			"app/usecase/usecase",
			"",
			errNotMoney,
		)
		return uuid.UUID{}, errNotMoney
	}

	fromAccount.Ballance = fromAccount.Ballance.Sub(transaction.Amount)
	toAccount.Ballance = toAccount.Ballance.Add(transaction.Amount)
	err = u.repo.TransferMoney(ctx, fromAccount, toAccount)
	if err != nil {
		return uuid.UUID{}, err
	}

	return u.repo.CreateTransaction(ctx, transaction)
}

func (u *usecase) GetTransactionsListByAccountID(ctx context.Context, accountID uint64) ([]models.Transaction, error) {
	return u.repo.GetTransactionsListByAccountID(ctx, accountID)
}
