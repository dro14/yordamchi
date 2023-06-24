package constants

import "time"

const (
	RetryAttempts = 10
)

const (
	RetryDelay      = 10 * time.Second
	RequestInterval = 2 * time.Second
)
