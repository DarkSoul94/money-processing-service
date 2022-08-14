package usecase

import "errors"

var (
	errInvalidTransactionType error = errors.New("invalid transaction type")
	errZeroAccountID          error = errors.New("account id is zero")
	errSameAccount            error = errors.New("the same account")
	errDifferentCurrencies    error = errors.New("accounts with different currencies")
	errNotMoney               error = errors.New("not enough money")
)
