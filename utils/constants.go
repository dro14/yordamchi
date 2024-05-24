package utils

import "time"

const (
	NumOfFreeReqs        = 50
	RetryAttempts        = 10
	RetryDelay           = 1000 * time.Millisecond
	ReqInterval          = 1000 * time.Millisecond
	NotificationInterval = 6 * time.Hour
	Delim                = "\n-\n-\n-\n-\n"
)
