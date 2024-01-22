package app

import (
	"github.com/gin-gonic/gin"
	"github.com/mauriciozanettisalomao/go-transaction-service/cmd/app/http/restapi"
)

// Endpoints returns the endpoints of the application
func Endpoints() *gin.Engine {

	r := gin.Default()

	r.POST("/v1/transactions", restapi.CreateTransaction)
	r.GET("/v1/transactions", restapi.ListTransactions)

	return r
}
