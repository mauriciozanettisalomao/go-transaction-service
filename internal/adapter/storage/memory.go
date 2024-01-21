package storage

import (
	"context"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/port"
)

// TransactionMemory is a memory storage for transactions created only for testing purposes
// It implements the TransactionWriter and TransactionRetriever interfaces
// With this implementation, we don't need to create a database, speeding up a containerized environment or
// any fancy solution available in the market for DBs in memory

var (
	transactions []domain.Transaction
	users        []domain.User
)

type transactionMemory struct {
}

func (m *transactionMemory) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	transactions = append(transactions, *transaction)
	return nil
}

func (m *transactionMemory) ListTransactions(ctx context.Context, limit int) ([]domain.Transaction, error) {
	return transactions, nil
}

func (m *transactionMemory) ValidateTransaction(ctx context.Context, transaction *domain.Transaction) error {

	return nil
}

// NewTransactionMemory creates an instance of a transaction memory
func NewTransactionMemory() port.TransactionHandler {
	return &transactionMemory{}
}

type userMemory struct {
}

func (m *userMemory) User(ctx context.Context, userID string) (*domain.User, error) {
	for _, user := range users {
		if user.ID == userID {
			return &user, nil
		}
	}
	return nil, nil
}

func (m *userMemory) CreateUser(ctx context.Context, user *domain.User) error {
	users = append(users, *user)
	return nil
}

// NewUserMemory creates an instance of a user memory
func NewUserMemory() port.UserRetriever {
	return &userMemory{}
}
