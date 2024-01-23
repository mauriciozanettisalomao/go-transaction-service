package domain

// Subscription represents a subscription to a topic
type Subscription struct {
	Topic    string `json:"-"`
	Protocol string `json:"protocol" binding:"required,oneof=email"`
	Endpoint string `json:"endpoint" binding:"required"`
}
