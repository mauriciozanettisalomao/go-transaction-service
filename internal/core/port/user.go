package port

import (
	"context"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
)

// UserRetriever defines the behavior of a user retriever
type UserRetriever interface {
	User(ctx context.Context, userID string) (*domain.User, error)
}

// UserWriter defines the behavior of a user writer
type UserWriter interface {
	CreateUser(ctx context.Context, user *domain.User) error
}
