package service

import (
	"context"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
)

// TransactionServiceMock implements the TransactionHandler interface for testing purposes
type TransactionServiceMock struct {
	CreateMock func(context.Context, *domain.Transaction) error
	ListMock   func(context.Context) ([]domain.Transaction, error)
}

// Create calls the CreateMock function
func (t *TransactionServiceMock) Create(ctx context.Context, transaction *domain.Transaction) error {
	return t.CreateMock(ctx, transaction)
}

// List calls the ListMock function
func (t *TransactionServiceMock) List(ctx context.Context) ([]domain.Transaction, error) {
	return t.ListMock(ctx)
}
