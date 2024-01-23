package service

import (
	"context"
	"os"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/port"
)

// SubscriptionHandler defines the behavior of a subscription handler
type SubscriptionHandler interface {
	Subscribe(context.Context, *domain.Subscription) error
}

// SubscriptionOptions helps to configure a subscription service
type SubscriptionOptions func(*subscriptorService)

// WithSubscriptor sets the subscriptor
func WithSubscriptor(svc port.Subscriptor) SubscriptionOptions {
	return func(ss *subscriptorService) {
		ss.adapter = svc
	}
}

type subscriptorService struct {
	adapter port.Subscriptor
}

// Subscribe subscribes to a topic
func (s *subscriptorService) Subscribe(ctx context.Context, subscription *domain.Subscription) error {
	if value, ok := os.LookupEnv("TRANSACTION_SNS_TOPIC_ARN"); ok {
		subscription.Topic = value
	}
	return s.adapter.Subscribe(ctx, subscription)
}

// NewSubscriptorService creates a new subscriptor service
func NewSubscriptorService(opts ...SubscriptionOptions) SubscriptionHandler {
	s := &subscriptorService{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
