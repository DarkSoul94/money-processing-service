package http

import (
	"context"
	"net/http"

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
