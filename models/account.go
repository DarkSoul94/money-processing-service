package models

import "github.com/shopspring/decimal"

type Account struct {
	Id       uint64
	Client   Client
	Currency Currency
	Ballance decimal.Decimal
}
