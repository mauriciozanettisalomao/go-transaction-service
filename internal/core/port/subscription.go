package port

import (
	"context"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
)

// Subscriptor defines the behavior of a subscriptor
type Subscriptor interface {
	Subscribe(ctx context.Context, subscription *domain.Subscription) error
}
