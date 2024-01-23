package restapi

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/adapter/notification"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/adapter/storage"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/port"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/service"
)

const (
	xIdempotencyKey = "x-idempotency-key"
)

type transactionAPI struct {
	svcTransaction  service.TransactionHandler
	svcSubscription service.SubscriptionHandler
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
//	@Security		X-API-Key
func (t *transactionAPI) CreateTransaction(ctx *gin.Context) {
	var req domain.Transaction
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	req.SetIdempontencyKey(ctx.GetHeader(xIdempotencyKey))
	if t.svcTransaction == nil {
		t.svcTransaction = service.NewTransactionHandler(
			service.WithTransactionWriter(storageLayerByEnv()),
			service.WithTransactionRetriever(storageLayerByEnv()),
		)
	}
	err := t.svcTransaction.Create(ctx, &req)
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
//	@Success		200						{object}	responseTransaction			"Successful operation"
//	@Failure		403						{object}	errorResponse			"Forbidden error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/v1/transactions [get]
//	@Security		X-API-Key
func (t *transactionAPI) ListTransactions(ctx *gin.Context) {

	fmt.Println(" limit", ctx.Query("limit"))
	limit := 10 // default limit
	limitParam, ok := ctx.GetQuery("limit")
	if ok {
		limitConv, err := strconv.Atoi(limitParam)
		if err != nil {
			slog.Error("error converting limit to int",
				"err", err,
				"limit", limitParam,
			)
			handleError(ctx, err)
		}
		limit = limitConv
	}

	response, next, err := service.NewTransactionHandler(
		service.WithTransactionRetriever(storageLayerByEnv()),
		service.WithNext(ctx.Query("next")),
		service.WithLimit(limit),
	).List(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, response, newMeta(limit, next))
}

// simple way to switch between memory and dynamodb
func storageLayerByEnv() port.TransactionHandler {
	if value, ok := os.LookupEnv("USE_DYNAMODB"); ok {
		if value == "true" {
			return storage.NewDynamoDB()
		}
	}
	return storage.NewTransactionMemory()
}

// SubscribeListenTransactions godoc
//
//	@Summary		Subscribe to listen the the new transactions
//	@Description	Subscribe to be notified when a new transaction is created
//	@Tags			transactions
//	@Accept			json
//	@Produce		json
//	@Param			Transaction				body		domain.Subscription	true	"Create Transaction request"
//	@Success		201						{object}	domain.Subscription			"Subscription created"
//	@Failure		403						{object}	errorResponse			"Forbidden error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/v1/transactions/subscribe [post]
//	@Security		X-API-Key
func (t *transactionAPI) SubscribeListenTransactions(ctx *gin.Context) {

	var req domain.Subscription
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	if t.svcSubscription == nil {
		t.svcSubscription = service.NewSubscriptorService(
			service.WithSubscriptor(notificationLayerByEnv()),
		)
	}
	err := t.svcSubscription.Subscribe(ctx, &req)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleCreatedSuccess(ctx, req)
}

// simple way to switch between memory and dynamodb
func notificationLayerByEnv() port.Subscriptor {
	if value, ok := os.LookupEnv("USE_SNS"); ok {
		if value == "true" {
			return notification.NewSubscriptionSnsTopic()
		}
	}
	return notification.NewSubscriptionMemory()
}

// NewTransactionAPI creates an instance of a transaction API
func NewTransactionAPI() *transactionAPI {
	return &transactionAPI{}
}
