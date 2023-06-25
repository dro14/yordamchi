package constants

import "time"

const (
	RetryAttempts = 10
)

const (
	RetryDelay      = 1 * time.Second
	RequestInterval = 2 * time.Second
)
