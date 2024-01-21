package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mauriciozanettisalomao/go-transaction-service/cmd/app/http/restapi"
)

func main() {

	r := gin.Default()

	r.POST("/transactions", restapi.CreateTransaction)
	r.GET("/transactions", restapi.ListTransactions)

	// Start server
	r.Run()
}
