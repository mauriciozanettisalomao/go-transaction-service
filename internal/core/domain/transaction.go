package domain

// Transaction represents a transaction domain model
// holds the business logic and functions to be used by the application
// to validate and manipulate the data
type Transaction struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Origin   string  `json:"origin"`
	User     struct {
		ID string `json:"id"`
	} `json:"user"`
	OperationType   string `json:"operationType"`
	CreatedAt       string `json:"createdAt"`
	IdempontencyKey string `json:"-"`
}

func (t *Transaction) Validate() error {

	// go routines to validate the fields

	return nil
}

func (t *Transaction) ValidateIdempotency(key string) error {
	if key == t.IdempontencyKey {
		return ErrDataAlreadyExists
	}
	return nil
}
