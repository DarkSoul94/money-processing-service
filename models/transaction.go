package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	Id     uuid.UUID
	Type   TransactionType
	From   Account
	To     Account
	Amount decimal.Decimal
}

type TransactionType int

const (
	Deposit TransactionType = iota
	Withdraw
	Transfer
)
