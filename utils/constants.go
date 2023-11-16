package utils

import "time"

const (
	NumOfFreeRequests = 5
	RetryAttempts     = 10
	RetryDelay        = 1000 * time.Millisecond
	RequestInterval   = 500 * time.Millisecond
)
