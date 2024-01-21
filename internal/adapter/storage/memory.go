package storage

import (
	"context"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
)

// TransactionMemory is a memory storage for transactions created only for testing purposes
// It implements the TransactionWriter and TransactionRetriever interfaces
// With this implementation, we don't need to create a database, speeding up a containerized environment or
// any fancy solution available in the market for DBs in memory

var (
	transactions []domain.Transaction
)

type memory struct {
}

func (m *memory) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	transactions = append(transactions, *transaction)
	return nil
}

func (m *memory) ListTransactions(ctx context.Context, limit int) ([]domain.Transaction, error) {
	return transactions, nil
}

func (m *memory) ValidateTransaction(ctx context.Context, transaction *domain.Transaction) error {

	for _, t := range transactions {
		if t.IdempontencyKey == transaction.IdempontencyKey {
			return t.ValidateIdempotency(transaction.IdempontencyKey)
		}
	}
	return nil
}

// NewTransactionMemory creates an instance of a transaction memory
func NewTransactionMemory() *memory {
	return &memory{}
}
