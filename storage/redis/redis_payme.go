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
	case "premium", "gpt-4":
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
		client.Set(ctx, "premium:"+id(ctx), expirationDate, time.Until(expiration))
	case "unlimited":
		expiration := time.Now()
		switch subscription {
		case "weekly":
			expiration = expiration.AddDate(0, 0, 7)
		case "monthly":
			expiration = expiration.AddDate(0, 1, 0)
		default:
			return fmt.Errorf("invalid subscription: %s", order)
		}
		expirationDate := expiration.Format("02.01.2006 15:04:05")
		client.Set(ctx, "unlimited:"+id(ctx), expirationDate, time.Until(expiration))
	case "images", "dall-e-3":
		purchased, err := strconv.Atoi(subscription)
		if err != nil {
			return fmt.Errorf("invalid number of purchased images: %s", subscription)
		}
		available, _ := client.Get(ctx, "images:"+id(ctx)).Int()
		client.Set(ctx, "images:"+id(ctx), available+purchased, 0)
	default:
		return fmt.Errorf("invalid order type: %s", order)
	}

	return nil
}
