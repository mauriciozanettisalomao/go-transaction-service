package app

import (
	"github.com/gin-gonic/gin"
	"github.com/mauriciozanettisalomao/go-transaction-service/cmd/app/http/restapi"
)

// Endpoints returns the endpoints of the application
func Endpoints() *gin.Engine {

	r := gin.Default()

	api := restapi.NewTransactionAPI()

	r.POST("/v1/transactions", api.CreateTransaction)
	r.GET("/v1/transactions", api.ListTransactions)

	return r
}
