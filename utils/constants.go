package utils

import "time"

const (
	NumOfFreeReqs = 5
	RetryAttempts = 10
	RetryDelay    = 1000 * time.Millisecond
	ReqInterval   = 2000 * time.Millisecond
	Delim         = "\n-\n-\n-\n-\n"
)
