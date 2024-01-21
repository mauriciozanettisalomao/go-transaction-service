package domain

// Transaction represents a transaction domain model
// holds the business logic and functions to be used by the application
// to validate and manipulate the data
type Transaction struct {
	ID            string  `json:"id"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Origin        string  `json:"origin"`
	UserID        string  `json:"user_id"`
	OperationType string  `json:"operation_type"`
	CreatedAt     string  `json:"created_at"`
}
