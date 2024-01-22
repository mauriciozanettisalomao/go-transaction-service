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
	transactionIdempontencyTable = "transaction_idempontency"
)

type dynamoDB struct {
	svc *dynamodb.DynamoDB
}

func (d *dynamoDB) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {

	if err := d.transaction(ctx, transaction); err != nil {
		return err
	}

	if err := d.idempotentTransaction(ctx, transaction.GetIdempontencyKey()); err != nil {
		return err
	}

	return nil
}

func (d *dynamoDB) transaction(ctx context.Context, transaction *domain.Transaction) error {

	type Transaction struct {
		ID            string
		UserID        string
		Currency      string
		Origin        string
		OperationType string
		CreatedAt     string
		Amount        float64
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

func (d *dynamoDB) idempotentTransaction(ctx context.Context, idempotencyKey string) error {

	type TransactionIdempotency struct {
		IdempontencyKey string `json:"idempontency_key"`
		CreatedAt       string `json:"created_at"`
	}

	ti := TransactionIdempotency{
		IdempontencyKey: idempotencyKey,
		CreatedAt:       time.Now().UTC().Format("2006-01-02 15:04:05"),
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
			"idempontency_key": {
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

	if response.Item["idempontency_key"] != nil {
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
