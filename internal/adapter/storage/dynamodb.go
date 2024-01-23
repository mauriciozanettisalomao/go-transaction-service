package storage

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

	if err := d.transaction(transaction); err != nil {
		return err
	}

	if err := d.idempotentTransaction(transaction.GetIdempontencyKey(), transaction.ID); err != nil {
		return err
	}

	return nil
}

func (d *dynamoDB) transaction(transaction *domain.Transaction) error {

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

func (d *dynamoDB) idempotentTransaction(idempotencyKey, transacionID string) error {

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

func (d *dynamoDB) ListTransactions(ctx context.Context, limit int, next string) ([]domain.Transaction, error) {

	input := &dynamodb.ScanInput{
		TableName: aws.String(transactionTableName),
		Limit:     aws.Int64(int64(limit)),
	}
	if next != "" {
		nextID, err := d.deserializeNextID(next)
		if err != nil {
			slog.Error("error deserializing the next id",
				"err", err,
				"next", next,
			)
			return nil, err
		}
		input.ExclusiveStartKey = nextID
	}

	output, err := d.svc.Scan(input)
	if err != nil {
		slog.Error("error getting the transaction by idempotency id from dynamodb",
			"err", err,
			"limit", limit,
		)
		return nil, err
	}

	nextID := ""
	if output.LastEvaluatedKey != nil {
		next, errNextId := d.nextID(*output)
		if errNextId != nil {
			return nil, errNextId
		}
		nextID = next
	}

	var transactions []domain.Transaction
	for _, i := range output.Items {
		var t domain.Transaction
		err = dynamodbattribute.UnmarshalMap(i, &t)
		if err != nil {
			slog.Error("error unmarshalling the transaction",
				"err", err,
				"item", i,
			)
			return nil, err
		}
		t.SetNext(nextID)
		t.User.ID = *i["userId"].S
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (d *dynamoDB) nextID(result dynamodb.ScanOutput) (string, error) {

	lekOutPut := ""
	if result.LastEvaluatedKey != nil {
		lek := map[string]interface{}{}
		err := dynamodbattribute.UnmarshalMap(result.LastEvaluatedKey, &lek)
		if err != nil {
			slog.Error("error unmarshalling the last evaluated key",
				"err", err,
				"lastEvaluatedKey", result.LastEvaluatedKey,
			)
			return lekOutPut, err
		}
		lastKey, err := json.Marshal(lek)
		if err != nil {
			slog.Error("error marshalling the last evaluated key",
				"err", err,
				"lastEvaluatedKey", result.LastEvaluatedKey,
			)
			return lekOutPut, err
		}
		lekOutPut = base64.StdEncoding.EncodeToString(lastKey)
	}

	return lekOutPut, nil

}

func (d *dynamoDB) deserializeNextID(input string) (map[string]*dynamodb.AttributeValue, error) {
	bytesJSON, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		slog.Error("error decoding the next id",
			"err", err,
			"input", input,
		)
		return nil, err
	}
	outputJSON := map[string]interface{}{}
	err = json.Unmarshal(bytesJSON, &outputJSON)
	if err != nil {
		slog.Error("error unmarshalling the next id",
			"err", err,
			"json", string(bytesJSON),
		)
		return nil, err
	}

	return dynamodbattribute.MarshalMap(outputJSON)
}

// NewDynamoDB returns a new DynamoDB instance
func NewDynamoDB() port.TransactionHandler {
	return &dynamoDB{
		svc: dynamodb.New(awsAdapter.Session()),
	}
}
