package port

import (
	"context"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
)

// UserRetriever defines the behavior of a user retriever
type UserRetriever interface {
	Get(ctx context.Context) (*domain.User, error)
}
