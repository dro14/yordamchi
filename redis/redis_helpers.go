package redis

import "time"

func untilMidnight() time.Duration {
	t := time.Now().AddDate(0, 0, 1)
	return time.Until(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local))
}

func midnight() string {
	t := time.Now().AddDate(0, 0, 1)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Format("15:04:05 02.01.2006")
}
