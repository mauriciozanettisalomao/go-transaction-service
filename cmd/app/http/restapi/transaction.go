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

// CreateTransaction godoc
//
//	@Summary		Create a new transaction
//	@Description	Create transactions made by a certain user
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			X-Idempotency-Key		header		string				true	"it helps you retry requests safely without accidentally doing the same thing twice. When making or changing an object, use an idempotency key."
//	@Param			Transaction				body		domain.Transaction	true	"Create Transaction request"
//	@Success		201						{object}	domain.Transaction		"Transaction created"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		403						{object}	errorResponse			"Forbidden error"
//	@Failure		404						{object}	errorResponse			"Data not found error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/v1/transactions [post]
//	@Security		ApiKeyAuth
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

// ListTransactions godoc
//
//	@Summary		Create a new transaction
//	@Description	Create transactions made by a certain user
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			limit					query		string				true	"The maximum number of records to return per page."
//	@Success		200						{object}	response				"Successful operation"
//	@Failure		403						{object}	errorResponse			"Forbidden error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/v1/transactions [get]
//	@Security		ApiKeyAuth
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
