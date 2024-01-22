package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionSetCreatedAt(t *testing.T) {

	tests := []struct {
		name     string
		date     string
		input    Transaction
		expected Transaction
	}{
		{
			name: "valid scenario",
			input: Transaction{
				ID: "999",
			},
			date: "2020-01-01T00:00:00Z",
			expected: Transaction{
				ID:        "999",
				CreatedAt: "2020-01-01T00:00:00Z",
			},
		},
	}

	assertion := assert.New(t)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.SetCreatedAt(tc.date)
			assertion.Equal(tc.expected, tc.input)
		})
	}
}

func TestTransactionSetID(t *testing.T) {

	tests := []struct {
		name     string
		id       string
		input    Transaction
		expected Transaction
	}{
		{
			name:  "valid scenario",
			id:    "abc",
			input: Transaction{},
			expected: Transaction{
				ID: "abc",
			},
		},
	}

	assertion := assert.New(t)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.SetID(tc.id)
			assertion.Equal(tc.expected, tc.input)
		})
	}
}

func TestTransactionIdempotencyID(t *testing.T) {

	tests := []struct {
		name            string
		idempontencyKey string
		input           Transaction
		expected        Transaction
	}{
		{
			name:            "valid scenario",
			idempontencyKey: "abc",
			input: Transaction{
				ID: "999",
			},
			expected: Transaction{
				idempontencyKey: "abc",
				ID:              "999",
			},
		},
	}

	assertion := assert.New(t)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.input.SetIdempontencyKey(tc.idempontencyKey)
			assertion.Equal(tc.expected.idempontencyKey, tc.input.GetIdempontencyKey())
		})
	}
}

func TestTransactionValidateIdempotency(t *testing.T) {

	tests := []struct {
		name       string
		validation error
		expected   error
		input      Transaction
	}{
		{
			name:       "with error",
			validation: errors.New("transaction already exists"),
			input: Transaction{
				ID: "999",
			},
			expected: ErrDataAlreadyExists,
		},
		{
			name:       "no error",
			validation: nil,
			input: Transaction{
				ID: "999",
			},
			expected: nil,
		},
	}

	assertion := assert.New(t)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ValidateIdempotency(tc.validation)
			assertion.Equal(tc.expected, result)
		})
	}
}

func TestTransactionValidateUserID(t *testing.T) {

	tests := []struct {
		name       string
		validation string
		expected   error
		input      Transaction
	}{
		{
			name:       "user does not exist",
			validation: "",
			input: Transaction{
				ID: "999",
				User: struct {
					ID string "json:\"id\"  binding:\"required\""
				}{
					ID: "user-id-abc",
				},
			},
			expected: ErrDataNotFound,
		},
		{
			name:       "user exists",
			validation: "user-id-abc",
			input: Transaction{
				ID: "999",
				User: struct {
					ID string "json:\"id\"  binding:\"required\""
				}{
					ID: "user-id-abc",
				},
			},
			expected: nil,
		},
	}

	assertion := assert.New(t)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ValidateUserID(tc.validation)
			assertion.Equal(tc.expected, result)
		})
	}
}

func TestTransactionAmount(t *testing.T) {

	tests := []struct {
		name       string
		validation string
		expected   error
		input      Transaction
	}{
		{
			name: "valid amount",
			input: Transaction{
				ID:     "999",
				Amount: 999.99,
			},
			expected: nil,
		},
		{
			name: "invalid amount",
			input: Transaction{
				ID:     "-999",
				Amount: -999.99,
			},
			expected: ErrInvalidAmount,
		},
	}

	assertion := assert.New(t)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.ValidateAmount()
			assertion.Equal(tc.expected, result)
		})
	}
}
