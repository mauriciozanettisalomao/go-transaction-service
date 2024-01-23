package notification

import (
	"context"
	"log/slog"

	awsAdapter "github.com/mauriciozanettisalomao/go-transaction-service/internal/adapter/aws"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/port"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type snsTopic struct {
	svc *sns.SNS
}

func (s *snsTopic) Subscribe(ctx context.Context, subscription *domain.Subscription) error {

	_, err := s.svc.Subscribe(&sns.SubscribeInput{
		Endpoint:              &subscription.Endpoint,
		Protocol:              &subscription.Protocol,
		ReturnSubscriptionArn: aws.Bool(true),
		TopicArn:              &subscription.Topic,
	})
	if err != nil {
		slog.Error("error subscribing to topic",
			"error", err,
			"topic", subscription.Topic,
			"endpoint", subscription.Endpoint,
			"protocol", subscription.Protocol,
		)
		return err
	}

	return nil
}

// NewSubscriptionSnsTopic creates an instance of a subscription sns topic
func NewSubscriptionSnsTopic() port.Subscriptor {
	return &snsTopic{
		svc: sns.New(awsAdapter.Session()),
	}
}
