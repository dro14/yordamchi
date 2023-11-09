package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/go-redis/redis/v8"
)

func isGPT4(ctx context.Context) (bool, error) {
	_, err := Client.Get(ctx, "model:"+id(ctx)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func isPremium(ctx context.Context) (bool, error) {
	_, err := Client.Get(ctx, "premium:"+id(ctx)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func isFree(ctx context.Context) (bool, error) {
	requests, err := Client.Get(ctx, "free:"+id(ctx)).Int()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			Client.Set(ctx, "free:"+id(ctx), constants.NumOfFreeRequests, untilMidnight())
			return true, nil
		}
		return false, err
	}
	if requests > 0 && requests <= constants.NumOfFreeRequests {
		return true, nil
	} else if requests == 0 {
		return false, nil
	} else {
		return false, fmt.Errorf("invalid number of requests for %s: %d", id(ctx), requests)
	}
}
