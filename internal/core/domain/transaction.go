package domain

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Transaction represents a transaction domain model
// holds the business logic and functions to be used by the application
// to validate and manipulate the data
type Transaction struct {
	User struct {
		ID string `json:"id"  binding:"required"`
	} `json:"user"  binding:"required"`
	ID              string `json:"id"`
	Currency        string `json:"currency"  binding:"required"`
	Origin          string `json:"origin"  binding:"required"`
	OperationType   string `json:"operationType"  binding:"required,oneof=debit credit"`
	CreatedAt       string `json:"createdAt"`
	idempontencyKey string
	Amount          float64 `json:"amount"  binding:"required"`
}

// Build sets default values for a transaction
func (t *Transaction) Build() *Transaction {
	t.SetID(uuid.New().String())
	t.SetCreatedAt(time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	return t
}

// SetCreatedAt sets the created at field for a transaction
func (t *Transaction) SetCreatedAt(date string) {
	t.CreatedAt = date
}

// SetID sets the id for a transaction
func (t *Transaction) SetID(id string) {
	t.ID = id
}

// SetIdempontencyKey sets the idempontency key for a transaction
func (t *Transaction) SetIdempontencyKey(key string) {
	t.idempontencyKey = key
}

// GetIdempontencyKey gets the idempontency key for a transaction
func (t *Transaction) GetIdempontencyKey() string {
	return t.idempontencyKey
}

// ValidateIdempotency validates if a transaction is idempotent
// the logic is just to handle which error should be returned
func (t *Transaction) ValidateIdempotency(err error) error {
	if err != nil {
		slog.Error("transaction already exists",
			"err", ErrDataAlreadyExists,
			"idempontencyKey", t.idempontencyKey,
		)
		return ErrDataAlreadyExists
	}
	return nil
}

// ValidateUserID validates if a transaction has a valid user id
// the logic is just to handle which error should be returned
func (t *Transaction) ValidateUserID(key string) error {
	if key != t.User.ID {
		slog.Error("user does not exist",
			"err", ErrDataNotFound,
			"userId", t.User.ID,
		)
		return ErrDataNotFound
	}
	return nil
}

// ValidateAmount validates if a transaction has a valid amount
func (t *Transaction) ValidateAmount() error {
	if t.Amount <= 0 {
		slog.Error("amount is invalid",
			"err", ErrInvalidAmount,
			"amount", t.Amount,
		)
		return ErrInvalidAmount
	}
	return nil
}
