package service

import (
	"context"
	"time"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/adapter/storage"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/port"

	"github.com/google/uuid"
)

// TransactionHandler defines the behavior of a transaction handler
type TransactionHandler interface {
	Create(context.Context, *domain.Transaction) error
	List(context.Context) ([]domain.Transaction, error)
}

type transactionService struct {
	transactionRetriever port.TransactionRetriever
	transactionWriter    port.TransactionWriter
	limit                int
}

func (ts *transactionService) Create(ctx context.Context, transaction *domain.Transaction) error {

	transaction.CreatedAt = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	transaction.ID = uuid.New().String()

	// validations := []func() error{
	// 	func() error {
	// 		return ts.transactionRetriever.Validate(ctx, transaction)
	// 	},
	// }

	// var wg sync.WaitGroup
	// wg.Add(len(validations))

	// for _, validation := range validations {
	// 	go func(v func() error) {
	// 		defer wg.Done()
	// 		v()
	// 	}(validation)
	// }
	// wg.Wait()

	return ts.transactionWriter.CreateTransaction(ctx, transaction)
}

func (ts *transactionService) List(ctx context.Context) ([]domain.Transaction, error) {
	ts.limit = 5
	return ts.transactionRetriever.ListTransactions(ctx, ts.limit)
}

// NewTransactionHandler creates an instance new transaction handler
func NewTransactionHandler() TransactionHandler {
	return &transactionService{
		transactionRetriever: storage.NewTransactionMemory(),
		transactionWriter:    storage.NewTransactionMemory(),
	}
}
