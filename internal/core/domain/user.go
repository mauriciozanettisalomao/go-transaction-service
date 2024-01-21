package domain

// User represents a user domain model holds the business logic and functions to
// be used by the application to validate and manipulate the data
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}
