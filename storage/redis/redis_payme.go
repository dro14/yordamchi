package redis

import (
	"context"
	"time"
)

func (r *Redis) PerformTransaction(ctx context.Context, amount int, Type string) {
	if Type == "gpt-4" {
		tokens, _ := r.client.Get(ctx, "gpt-4:"+id(ctx)).Int()
		tokens += amount / 100
		r.client.Set(ctx, "gpt-4:"+id(ctx), tokens, 0)
	} else {
		var expiration time.Time
		switch Type {
		case "weekly":
			expiration = time.Now().AddDate(0, 0, 7)
		case "monthly":
			expiration = time.Now().AddDate(0, 1, 0)
		}
		expirationDate := expiration.Format("15:04:05 02.01.2006")
		r.client.Set(ctx, "premium:"+id(ctx), expirationDate, time.Until(expiration))
	}
}
