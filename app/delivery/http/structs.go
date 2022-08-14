package http

import (
	"github.com/DarkSoul94/money-processing-service/consts"
	"github.com/DarkSoul94/money-processing-service/models"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type newClient struct {
	Name string `json:"name"`
}

func (h *Handler) toModelClient(client newClient) models.Client {
	return models.Client{
		Name: client.Name,
	}
}

type outClient struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) toOutClient(client models.Client) outClient {
	return outClient{
		Id:   client.Id,
		Name: client.Name,
	}
}

type newAccount struct {
	ClientID   uint64 `json:"client_id"`
	CurrencyID uint   `json:"currency_id"`
}

type outCurrency struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) toOutCurrency(currency models.Currency) outCurrency {
	return outCurrency{
		Id:   currency.Id,
		Name: currency.Name,
	}
}

type outAccount struct {
	Id       uint64          `json:"id"`
	Currency outCurrency     `json:"currency"`
	Ballance decimal.Decimal `json:"ballance"`
}

func (h *Handler) toOutAccount(account models.Account) outAccount {
	return outAccount{
		Id:       account.Id,
		Currency: h.toOutCurrency(account.Currency),
		Ballance: account.Ballance,
	}
}

type newTransaction struct {
	Type          uint            `json:"type"`
	FromAccountID uint64          `json:"from_account_id,omitempty"`
	ToAccountID   uint64          `json:"to_account_id,omitempty"`
	Amount        decimal.Decimal `json:"amount"`
}

func (h *Handler) toModelTransction(transaction newTransaction) models.Transaction {
	return models.Transaction{
		Type: models.TransactionType{Id: transaction.Type},
		From: models.Account{
			Id: transaction.FromAccountID,
		},
		To: models.Account{
			Id: transaction.ToAccountID,
		},
		Amount: transaction.Amount,
	}
}

type outTransaction struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Type      string `json:"type"`
	From      uint64 `json:"from"`
	To        uint64 `json:"to"`
	Amount    string `json:"amount"`
}

func (h *Handler) toOutTransaction(mTransaction models.Transaction) outTransaction {
	transaction := outTransaction{
		Id:        mTransaction.Id.String(),
		CreatedAt: mTransaction.CreatedAt.Format(consts.OutTransactionTime),
		From:      mTransaction.From.Id,
		To:        mTransaction.To.Id,
		Amount:    mTransaction.Amount.String(),
	}

	switch mTransaction.Type.Id {
	case models.Deposit.Id:
		transaction.Type = models.Deposit.Name
	case models.Withdraw.Id:
		transaction.Type = models.Withdraw.Name
	case models.Transfer.Id:
		transaction.Type = models.Transfer.Name
	}

	return transaction
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) newErrorResponse(c *gin.Context, statusCode int, err error) {
	c.AbortWithStatusJSON(statusCode, errorResponse{err.Error()})
}
