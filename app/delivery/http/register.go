package http

import (
	"github.com/DarkSoul94/money-processing-service/app"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc app.Usecase) {
	h := NewHandler(uc)

	clientApi := router.Group("/client")
	{
		clientApi.POST("", h.CreateClient)
		clientApi.GET(":id", h.GetClientByID)
	}

	accountApi := router.Group("/account")
	{
		accountApi.POST("", h.CreateAccount)
		accountApi.GET(":id", h.GetAccountByID)
	}

	transactionApi := router.Group("/transaction")
	{
		transactionApi.POST("", h.CreateTransaction)
		transactionApi.GET(":id", h.GetTransactionsListByAccountID)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
