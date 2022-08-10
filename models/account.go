package models

import "google.golang.org/genproto/googleapis/type/decimal"

type Account struct {
	Id       uint64
	Client   Client
	Currency Currency
	Ballance decimal.Decimal
}
