package utils

import "time"

const (
	NumOfFreeReqs        = 10
	RetryAttempts        = 10
	RetryDelay           = 1000 * time.Millisecond
	ReqInterval          = 1000 * time.Millisecond
	NotificationInterval = 6 * time.Hour
	AdminUserID          = 1331278972
	Delim                = "\n-\n-\n-\n-\n"
)
