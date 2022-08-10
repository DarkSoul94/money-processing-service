package models

import "google.golang.org/genproto/googleapis/type/decimal"

type Transaction struct {
	Id         uint64
	Type       TransactionType
	From       Account
	To         Account
	MoneyValue decimal.Decimal
}

type TransactionType int

const (
	Deposit TransactionType = iota
	Withdraw
	Transfer
)
