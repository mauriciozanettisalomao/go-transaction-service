package domain

import "errors"

var (
	// ErrDataNotFound is the error returned when the data is not found
	ErrDataNotFound = errors.New("data not found")
	// ErrInvalidData is the error returned when the data is invalid
	ErrInvalidData = errors.New("invalid data")
	// ErrInvalidPagination is the error returned when the pagination is invalid
	ErrInvalidPagination = errors.New("invalid pagination")
	// ErrInvalidUser is the error returned when the user is invalid
	ErrInvalidUser = errors.New("invalid user")
	// ErrInvalidAmount is the error returned when the amount is invalid
	ErrInvalidAmount = errors.New("invalid amount")
	// ErrInvalidCurrency is the error returned when the currency is invalid
	ErrDataAlreadyExists = errors.New("data already exists")
)
