package notification

import (
	"context"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/port"
)

// subscriptionMemory is a memory object for subscriptions
// It implements the Subscriptor interface
// With this implementation, we don't need to create a database, speeding up a containerized environment or
// any fancy solution available in the market for DBs in memory

var (
	subscriptions []domain.Subscription
)

type subscriptionMemory struct {
}

func (sm *subscriptionMemory) Subscribe(ctx context.Context, subscription *domain.Subscription) error {
	subscriptions = append(subscriptions, *subscription)
	return nil
}

// NewSubscriptionMemory creates an instance of a subscription memory
func NewSubscriptionMemory() port.Subscriptor {
	return &subscriptionMemory{}
}
