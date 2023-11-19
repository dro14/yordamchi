package redis

import (
	"context"
	"fmt"
	"time"
)

func (r *Redis) PerformTransaction(ctx context.Context, Type string) error {
	var expiration time.Time
	switch Type {
	case "daily:gpt-4":
		expiration = time.Now().AddDate(0, 0, 1)
	case "weekly:gpt-4":
		expiration = time.Now().AddDate(0, 0, 7)
	case "monthly:gpt-4":
		expiration = time.Now().AddDate(0, 1, 0)
	default:
		return fmt.Errorf("invalid transaction type: %s", Type)
	}
	expirationDate := expiration.Format("02.01.2006 15:04:05")
	r.client.Set(ctx, "premium:"+id(ctx), expirationDate, time.Until(expiration))
	return nil
}
