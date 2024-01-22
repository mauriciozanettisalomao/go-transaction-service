package port

import (
	"context"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
)

// UserRetriever defines the behavior of a user retriever
type UserRetriever interface {
	User(ctx context.Context, userID string) (*domain.User, error)
}
