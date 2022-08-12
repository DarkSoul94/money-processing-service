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

type TransactionType struct {
	Id   uint
	Name string
}

var (
	Deposit  TransactionType = TransactionType{Id: 0, Name: "Deposit"}
	Withdraw TransactionType = TransactionType{Id: 1, Name: "Withdraw"}
	Transfer TransactionType = TransactionType{Id: 2, Name: "Transfer"}
)
