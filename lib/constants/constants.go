package constants

import "time"

const (
	RetryAttempts = 10
)

const (
	RetryDelay      = 1000 * time.Millisecond
	RequestInterval = 2000 * time.Millisecond
)

var (
	MerchantID      string
	TestKey         string
	RealKey         string
	SubscriptionKey string
)
