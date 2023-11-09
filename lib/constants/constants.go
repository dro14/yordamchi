package constants

import "time"

const (
	RetryAttempts = 10
)

const (
	RetryDelay      = 1000 * time.Millisecond
	RequestInterval = 500 * time.Millisecond
)

const (
	NumOfFreeRequests = 5
)

var (
	MerchantID      string
	TestKey         string
	RealKey         string
	SubscriptionKey string
)
