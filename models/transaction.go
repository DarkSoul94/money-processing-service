package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	Id        uuid.UUID
	CreatedAt time.Time
	Type      TransactionType
	From      Account
	To        Account
	Amount    decimal.Decimal
}

type TransactionType int

const (
	Deposit TransactionType = iota
	Withdraw
	Transfer
)
