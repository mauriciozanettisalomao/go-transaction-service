package restapi

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/adapter/storage"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/port"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/service"
)

const (
	xIdempotencyKey = "x-idempotency-key"
)

type transactionAPI struct {
	svc service.TransactionHandler
}

// CreateTransaction defines the endpoint to create a transaction
func (t *transactionAPI) CreateTransaction(ctx *gin.Context) {
	var req domain.Transaction
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	req.SetIdempontencyKey(ctx.GetHeader(xIdempotencyKey))
	if t.svc == nil {
		t.svc = service.NewTransactionHandler(
			service.WithTransactionWriter(storageLayerByEnv()),
			service.WithTransactionRetriever(storageLayerByEnv()),
		)
	}
	err := t.svc.Create(ctx, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleCreatedSuccess(ctx, req)
}

// ListTransactions defines the endpoint to list transactions
// It accepts a query param limit to limit the number of transactions returned
func (t *transactionAPI) ListTransactions(ctx *gin.Context) {

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
		service.WithTransactionRetriever(storageLayerByEnv()),
		service.WithLimit(limit),
	).List(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, response, newMeta(limit))
}

// not a best way to check if it is running in a lambda environment
// it should be refactored to use a proper environment variable
func storageLayerByEnv() port.TransactionHandler {
	if _, ok := os.LookupEnv("AWS_LAMBDA_FUNCTION_NAME"); ok {
		return storage.NewDynamoDB()
	}
	return storage.NewTransactionMemory()
}

// NewTransactionAPI creates an instance of a transaction API
func NewTransactionAPI() *transactionAPI {
	return &transactionAPI{}
}
