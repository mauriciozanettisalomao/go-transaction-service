package app

import (
	"github.com/gin-gonic/gin"
	"github.com/mauriciozanettisalomao/go-transaction-service/cmd/app/http/middleware"
	"github.com/mauriciozanettisalomao/go-transaction-service/cmd/app/http/restapi"

	_ "github.com/mauriciozanettisalomao/go-transaction-service/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Endpoints returns the endpoints of the application
func Endpoints() *gin.Engine {

	r := gin.Default()

	r.Use(
		middleware.RequestIdMiddleware(),
		middleware.ExecutionTime(),
	)

	api := restapi.NewTransactionAPI()

	r.POST("/v1/transactions", api.CreateTransaction)
	r.POST("/v1/transactions/subscribe", api.SubscribeListenTransactions)
	r.GET("/v1/transactions", api.ListTransactions)

	r.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
