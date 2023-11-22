package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (r *Redis) PerformTransaction(ctx context.Context, order string) error {
	subscription, orderType, found := strings.Cut(order, ":")
	if !found {
		return fmt.Errorf("invalid order: %s", order)
	}

	switch orderType {
	case "gpt-4":
		expiration := time.Now()
		switch subscription {
		case "daily":
			expiration = expiration.AddDate(0, 0, 1)
		case "weekly":
			expiration = expiration.AddDate(0, 0, 7)
		case "monthly":
			expiration = expiration.AddDate(0, 1, 0)
		default:
			return fmt.Errorf("invalid subscription: %s", order)
		}
		expirationDate := expiration.Format("02.01.2006 15:04:05")
		r.client.Set(ctx, "premium:"+id(ctx), expirationDate, time.Until(expiration))
	case "dall-e-3":
		images, err := strconv.Atoi(subscription)
		if err != nil {
			return fmt.Errorf("invalid number of images to set: %s", subscription)
		}
		r.client.Set(ctx, "images:"+id(ctx), images, 0)
	default:
		return fmt.Errorf("invalid order type: %s", order)
	}

	return nil
}
