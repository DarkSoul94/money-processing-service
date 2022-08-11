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

func (h *Handler) CreateClient(c *gin.Context) {
	var client newClient

	if err := c.BindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"status": "error", "error": err.Error()})
		return
	}

	ctx, cansel := context.WithCancel(c)
	defer cansel()

	id, err := h.uc.CreateClient(ctx, h.toModelClient(client))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "client_id": id})
}

func (h *Handler) GetClientByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"status": "error", "error": "Invalid value in param 'id'"})
		return
	}

	ctx, cansel := context.WithCancel(c)
	defer cansel()

	mClient, idList, err := h.uc.GetClientByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "client": h.toOutClient(mClient), "accounts_id": idList})
}

func (h *Handler) CreateAccount(c *gin.Context) {
	var account newAccount

	if err := c.BindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"status": "error", "error": err.Error()})
		return
	}

	ctx, cansel := context.WithCancel(c)
	defer cansel()

	id, err := h.uc.CreateAccount(ctx, h.toModelAccount(account))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "account_id": id})
}

func (h *Handler) GetAccountByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"status": "error", "error": "Invalid value in param 'id'"})
		return
	}

	ctx, cansel := context.WithCancel(c)
	defer cansel()

	mAccount, err := h.uc.GetAccountByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "account": h.toOutAccount(mAccount)})
}

func (h *Handler) CreateTransaction(c *gin.Context) {
	var transaction newTransaction

	if err := c.BindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"status": "error", "error": err.Error()})
		return
	}

	ctx, cansel := context.WithCancel(c)
	defer cansel()

	id, err := h.uc.CreateTransaction(ctx, h.toModelTransction(transaction))
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "error", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "transaction_id": id})

}
