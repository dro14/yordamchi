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
	case "premium":
		expiration := time.Now()
		requests := ""
		switch subscription {
		case "daily":
			expiration = expiration.AddDate(0, 0, 1)
			requests = "50"
		case "weekly":
			expiration = expiration.AddDate(0, 0, 7)
			requests = "250"
		case "monthly":
			expiration = expiration.AddDate(0, 1, 0)
			requests = "750"
		default:
			return fmt.Errorf("invalid subscription: %s", order)
		}
		value := expiration.Format("02.01.2006 15:04:05") + "|" + requests
		client.Set(ctx, "premium:"+id(ctx), value, time.Until(expiration))
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
		value := expiration.Format("02.01.2006 15:04:05")
		client.Set(ctx, "unlimited:"+id(ctx), value, time.Until(expiration))
	case "images":
		images, err := strconv.Atoi(subscription)
		if err != nil {
			return fmt.Errorf("invalid number of images: %s", subscription)
		}
		client.IncrBy(ctx, "images:"+id(ctx), int64(images))
	default:
		return fmt.Errorf("invalid order type: %s", order)
	}

	return nil
}
