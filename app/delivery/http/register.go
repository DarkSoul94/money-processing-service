package http

import (
	"github.com/DarkSoul94/money-processing-service/app"
	"github.com/gin-gonic/gin"
)

// RegisterHTTPEndpoints ...
func RegisterHTTPEndpoints(router *gin.RouterGroup, uc app.Usecase) {
	h := NewHandler(uc)

	router.POST("/client", h.CreateClient)
	router.GET("/client/:id", h.GetClientByID)
	
	router.POST("/account", h.CreateAccount)
}
