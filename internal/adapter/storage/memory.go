package storage

import (
	"context"
	"errors"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/port"
)

// TransactionMemory is a memory storage for transactions created only for testing purposes
// It implements the TransactionWriter and TransactionRetriever interfaces
// With this implementation, we don't need to create a database, speeding up a containerized environment or
// any fancy solution available in the market for DBs in memory

var (
	transactions []domain.Transaction
)

type transactionMemory struct {
}

func (m *transactionMemory) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	transactions = append(transactions, *transaction)
	return nil
}

func (m *transactionMemory) ListTransactions(ctx context.Context, limit int, next string) ([]domain.Transaction, error) {
	return transactions, nil
}

func (m *transactionMemory) ValidateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	for _, t := range transactions {
		if t.GetIdempontencyKey() == transaction.GetIdempontencyKey() {
			return errors.New("transaction already exists")
		}
	}
	return nil
}

// NewTransactionMemory creates an instance of a transaction memory
func NewTransactionMemory() port.TransactionHandler {
	return &transactionMemory{}
}
