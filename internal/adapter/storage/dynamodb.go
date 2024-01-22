package storage

import (
	"context"
	"errors"
	"log/slog"
	"time"

	awsAdapter "github.com/mauriciozanettisalomao/go-transaction-service/internal/adapter/aws"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/port"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	transactionTableName         = "transaction"
	transactionIdempontencyTable = "transaction_idempotency"
)

type dynamoDB struct {
	svc *dynamodb.DynamoDB
}

func (d *dynamoDB) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {

	if err := d.transaction(ctx, transaction); err != nil {
		return err
	}

	if err := d.idempotentTransaction(ctx, transaction.GetIdempontencyKey(), transaction.ID); err != nil {
		return err
	}

	return nil
}

func (d *dynamoDB) transaction(ctx context.Context, transaction *domain.Transaction) error {

	type Transaction struct {
		ID            string  `json:"id"`
		UserID        string  `json:"userId"`
		Currency      string  `json:"currency"`
		Origin        string  `json:"origin"`
		OperationType string  `json:"operationType"`
		CreatedAt     string  `json:"createdAt"`
		Amount        float64 `json:"amount"`
	}

	t := Transaction{
		ID:            transaction.ID,
		UserID:        transaction.User.ID,
		Currency:      transaction.Currency,
		Origin:        transaction.Origin,
		OperationType: transaction.OperationType,
		Amount:        transaction.Amount,
		CreatedAt:     transaction.CreatedAt,
	}

	av, err := dynamodbattribute.MarshalMap(t)
	if err != nil {
		slog.Error("error marshalling the transaction",
			"err", err,
			"transaction", transaction,
		)
		return err
	}

	_, err = d.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(transactionTableName),
		Item:      av,
	})
	if err != nil {
		slog.Error("error inserting the transaction into dynamodb",
			"err", err,
			"transaction", transaction,
		)
		return err
	}

	return nil
}

func (d *dynamoDB) idempotentTransaction(ctx context.Context, idempotencyKey, transacionID string) error {

	if idempotencyKey == "" {
		return nil
	}

	type TransactionIdempotency struct {
		IdempotencyKey string `json:"idempotencyKey"`
		TransactionID  string `json:"transactionId"`
		CreatedAt      string `json:"createdAt"`
	}

	ti := TransactionIdempotency{
		IdempotencyKey: idempotencyKey,
		TransactionID:  transacionID,
		CreatedAt:      time.Now().UTC().Format("2006-01-02 15:04:05"),
	}

	av, err := dynamodbattribute.MarshalMap(ti)
	if err != nil {
		slog.Error("error marshalling the transaction idempotency",
			"err", err,
			"idempotencyKey", idempotencyKey,
		)
		return err
	}

	_, err = d.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(transactionIdempontencyTable),
		Item:      av,
	})
	if err != nil {
		slog.Error("error inserting the transaction into dynamodb",
			"err", err,
			"idempotencyKey", idempotencyKey,
		)
		return err
	}

	return nil
}

func (d *dynamoDB) ValidateTransaction(ctx context.Context, transaction *domain.Transaction) error {

	response, err := d.svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"idempotencyKey": {
				S: aws.String(transaction.GetIdempontencyKey()),
			},
		},
		TableName: aws.String(transactionIdempontencyTable),
	})
	if err != nil {
		slog.Error("error getting the transaction by idempotency id from dynamodb",
			"err", err,
			"idempotencyKey", transaction.GetIdempontencyKey(),
		)
		return err
	}

	if len(response.Item) > 0 {
		err := errors.New("transaction already exists")
		slog.Error("error getting the transaction by idempotency id from dynamodb",
			"err", err,
			"idempotencyKey", transaction.GetIdempontencyKey(),
		)
		return err
	}

	return nil
}

func (d *dynamoDB) ListTransactions(ctx context.Context, limit int) ([]domain.Transaction, error) {
	return nil, nil
}

// NewDynamoDB returns a new DynamoDB instance
func NewDynamoDB() port.TransactionHandler {
	return &dynamoDB{
		svc: dynamodb.New(awsAdapter.Session()),
	}
}
