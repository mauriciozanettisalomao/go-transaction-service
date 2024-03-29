package port

import (
	"context"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
)

// TransactionHandler defines the behavior of a transaction handler
type TransactionHandler interface {
	TransactionWriter
	TransactionRetriever
}

// TransactionRetriever defines the behavior of a transaction retriever
type TransactionRetriever interface {
	ListTransactions(ctx context.Context, limit int, next string) ([]domain.Transaction, error)
	ValidateTransaction(ctx context.Context, transaction *domain.Transaction) error
}

// TransactionWriter defines the behavior of a transaction writer
type TransactionWriter interface {
	CreateTransaction(ctx context.Context, transaction *domain.Transaction) error
}
