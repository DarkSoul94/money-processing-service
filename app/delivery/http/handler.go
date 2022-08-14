package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/DarkSoul94/money-processing-service/app"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc app.Usecase
}

func NewHandler(uc app.Usecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

// CreateClient godoc
// @Summary 			Create client
// @Description  	Create client
// @Accept       	json
// @Produce      	json
// @Param 				input body newClient true "Client name"
// @Success 			200 {integer} integer "Client ID"
// @Failure 			400 {object} errorResponse
// @Failure 			500 {object} errorResponse
// @Router 				/client [post]
func (h *Handler) CreateClient(c *gin.Context) {
	var client newClient

	if err := c.BindJSON(&client); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithCancel(c)
	defer cancel()

	id, err := h.uc.CreateClient(ctx, h.toModelClient(client))
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"client_id": id})
}

// GetClientByID 	godoc
// @Summary 			Get client
// @Description  	Get client by id
// @Accept       	json
// @Produce      	json
// @Param 				id path integer true "Client id"
// @Success 			200 {object} outClient
// @Success 			200 {array} integer
// @Failure 			400 {object} errorResponse
// @Failure 			500 {object} errorResponse
// @Router 				/client/{id} [get]
func (h *Handler) GetClientByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		h.newErrorResponse(c, http.StatusBadRequest, errInvalidID)
		return
	}

	ctx, cancel := context.WithCancel(c)
	defer cancel()

	mClient, idList, err := h.uc.GetClientByID(ctx, id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"client": h.toOutClient(mClient), "accounts_id": idList})
}

// CreateAccount godoc
// @Summary 			Create account
// @Description  	Create account
// @Accept       	json
// @Produce      	json
// @Param 				input body newAccount true "Client id"
// @Success 			200 {integer} integer "Account ID"
// @Failure 			400 {object} errorResponse
// @Failure 			500 {object} errorResponse
// @Router 				/account [post]
func (h *Handler) CreateAccount(c *gin.Context) {
	var account newAccount

	if err := c.BindJSON(&account); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithCancel(c)
	defer cancel()

	id, err := h.uc.CreateAccount(ctx, account.CurrencyID, account.ClientID)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"account_id": id})
}

// GetAccountByID godoc
// @Summary 			Get account
// @Description  	Get account by id
// @Accept       	json
// @Produce      	json
// @Param 				id path integer true "Account id"
// @Success 			200 {object} outAccount
// @Failure 			400 {object} errorResponse
// @Failure 			500 {object} errorResponse
// @Router 				/account/{id} [get]
func (h *Handler) GetAccountByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		h.newErrorResponse(c, http.StatusBadRequest, errInvalidID)
		return
	}

	ctx, cancel := context.WithCancel(c)
	defer cancel()

	mAccount, err := h.uc.GetAccountByID(ctx, id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"account": h.toOutAccount(mAccount)})
}

// CreateTransaction godoc
// @Summary 			Create transaction
// @Description  	Create transaction
// @Accept       	json
// @Produce      	json
// @Param 				input body newTransaction true "transaction type"
// @Success 			200 {integer} integer "Transaction ID"
// @Failure 			400 {object} errorResponse
// @Failure 			500 {object} errorResponse
// @Router 				/transaction [post]
func (h *Handler) CreateTransaction(c *gin.Context) {
	var transaction newTransaction

	if err := c.BindJSON(&transaction); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithCancel(c)
	defer cancel()

	id, err := h.uc.CreateTransaction(ctx, h.toModelTransction(transaction))
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"transaction_id": id})
}

// GetTransactionsListByAccountID godoc
// @Summary 											Get transactions list
// @Description  									Get transactions list by account id
// @Accept       									json
// @Produce      									json
// @Param 												id path integer true "Account id"
// @Success 											200 {array} outTransaction
// @Failure 			400 {object} errorResponse
// @Failure 			500 {object} errorResponse
// @Router 												/transaction/{id} [get]
func (h *Handler) GetTransactionsListByAccountID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		h.newErrorResponse(c, http.StatusBadRequest, errInvalidID)
		return
	}

	ctx, cancel := context.WithCancel(c)
	defer cancel()

	mTransactionsList, err := h.uc.GetTransactionsListByAccountID(ctx, id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	outTransactionsList := make([]outTransaction, 0)
	for _, mTransaction := range mTransactionsList {
		outTransactionsList = append(outTransactionsList, h.toOutTransaction(mTransaction))
	}

	c.JSON(http.StatusOK, map[string]interface{}{"transactions": outTransactionsList})
}
