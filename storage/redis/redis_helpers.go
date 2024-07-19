package redis

import (
	"context"
	"fmt"
	"time"
)

func id(ctx context.Context) string {
	return fmt.Sprintf("%d", ctx.Value("user_id").(int64))
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

func untilMidnight() time.Duration {
	t := time.Now().AddDate(0, 0, 1)
	return time.Until(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local))
}

func midnight() string {
	t := time.Now().AddDate(0, 0, 1)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Format("02.01.2006 15:04:05")
}
