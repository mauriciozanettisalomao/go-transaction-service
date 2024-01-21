package port

import (
	"context"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
)

// TransactionRetriever defines the behavior of a transaction retriever
type TransactionRetriever interface {
	List(ctx context.Context, limit uint64) ([]domain.Transaction, error)
}

// TransactionWriter defines the behavior of a transaction writer
type TransactionWriter interface {
	Create(ctx context.Context, transaction domain.Transaction) error
}
