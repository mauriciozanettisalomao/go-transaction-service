package restapi

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/adapter/storage"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/service"
)

const (
	xIdempotencyKey = "x-idempotency-key"
)

// CreateTransaction defines the endpoint to create a transaction
func CreateTransaction(ctx *gin.Context) {
	var req domain.Transaction
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	req.SetIdempontencyKey(ctx.GetHeader(xIdempotencyKey))
	err := service.NewTransactionHandler(
		service.WithTransactionWriter(storage.NewTransactionMemory()),
		service.WithTransactionRetriever(storage.NewTransactionMemory()),
	).Create(ctx, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleCreatedSuccess(ctx, req)
}

// ListTransactions defines the endpoint to list transactions
// It accepts a query param limit to limit the number of transactions returned
func ListTransactions(ctx *gin.Context) {

	fmt.Println(" limit", ctx.Query("limit"))
	limit := 10 // default limit
	limitParam, ok := ctx.GetQuery("limit")
	if ok {
		limitConv, err := strconv.Atoi(limitParam)
		if err != nil {
			fmt.Println(err)

		}
		limit = limitConv
	}

	response, err := service.NewTransactionHandler(
		service.WithTransactionRetriever(storage.NewTransactionMemory()),
		service.WithLimit(limit),
	).List(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, response, newMeta(limit))
}
