package service

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"

	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/domain"
	"github.com/mauriciozanettisalomao/go-transaction-service/internal/core/port"

	"golang.org/x/sync/errgroup"
)

// TransactionHandler defines the behavior of a transaction handler
type TransactionHandler interface {
	Create(context.Context, *domain.Transaction) error
	List(context.Context) ([]domain.Transaction, string, error)
}

type transactionService struct {
	transactionRetriever port.TransactionRetriever
	transactionWriter    port.TransactionWriter
	next                 string
	limit                int
}

// TransactionOptions helps to configure a transaction service
type TransactionOptions func(*transactionService)

// WithTransactionRetriever sets the transaction retriever
func WithTransactionRetriever(transactionRetriever port.TransactionRetriever) TransactionOptions {
	return func(ts *transactionService) {
		ts.transactionRetriever = transactionRetriever
	}
}

// WithTransactionWriter sets the transaction writer
func WithTransactionWriter(transactionWriter port.TransactionWriter) TransactionOptions {
	return func(ts *transactionService) {
		ts.transactionWriter = transactionWriter
	}
}

// WithNext sets the next token for a paginated response
func WithNext(next string) TransactionOptions {
	return func(ts *transactionService) {
		ts.next = next
	}
}

// WithLimit sets the limit for a paginated response
func WithLimit(limit int) TransactionOptions {
	return func(ts *transactionService) {
		ts.limit = limit
	}
}

func (ts *transactionService) Create(ctx context.Context, transaction *domain.Transaction) error {

	slog.Debug("creating transaction",
		"transaction", transaction,
		"requestID", ctx.Value("requestID"),
	)

	errs, ctx := errgroup.WithContext(ctx)

	// validations to be executed in parallel
	validations := []func() error{
		func() error {
			if transaction.GetIdempontencyKey() == "" {
				return nil
			}
			err := ts.transactionRetriever.ValidateTransaction(ctx, transaction)
			if err != nil {
				slog.Error("error validating transaction",
					"err", err,
					"transaction", transaction,
					"requestID", ctx.Value("requestID"),
				)
			}
			return transaction.ValidateIdempotency(err)
		},
		func() error {
			// once the solution to handle users is implemented, this validation should be implemented
			// Solutions: auth0, cognito, in-house solution, etc
			return nil
		},
		func() error {
			return transaction.ValidateAmount()
		},
	}

	validationChan := make(chan func() error, len(validations))
	// producer
	go func() {
		defer close(validationChan)
		for _, validation := range validations {
			validationChan <- validation
		}
	}()

	for i := 0; i < runtime.NumCPU(); i++ {
		idx := i
		// consumer
		errs.Go(func() (err error) {

			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("recovered for %v", r)
					slog.Error("error validating request",
						"err", err,
						"requestID", ctx.Value("requestID"),
					)
				}
			}()

			for f := range validationChan {

				slog.Debug("starting validation",
					"routine", idx,
					"requestID", ctx.Value("requestID"),
				)

				err = f()
				if err != nil {
					slog.Error("error validation request",
						"err", err,
						"routine", idx,
						"requestID", ctx.Value("requestID"),
					)
					break
				}
			}
			return err
		})

	}

	err := errs.Wait()
	if err != nil {
		slog.Error("request is not valid",
			"err", err,
			"requestID", ctx.Value("requestID"),
		)
		return err
	}

	return ts.transactionWriter.CreateTransaction(ctx, transaction.Build())
}

func (ts *transactionService) List(ctx context.Context) ([]domain.Transaction, string, error) {
	transactions, err := ts.transactionRetriever.ListTransactions(ctx, ts.limit, ts.next)
	if err != nil {
		slog.Error("error listing transactions",
			"err", err,
			"requestID", ctx.Value("requestID"),
		)
		return nil, "", err
	}

	nextKey := ""
	for i := range transactions {
		nextKey = transactions[i].GetNext()
	}

	return transactions, nextKey, nil
}

// NewTransactionHandler creates an instance new transaction handler
func NewTransactionHandler(opts ...TransactionOptions) TransactionHandler {
	ts := &transactionService{}
	for _, opt := range opts {
		opt(ts)
	}
	return ts
}
