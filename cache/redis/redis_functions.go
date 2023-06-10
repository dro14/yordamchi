package redis

import (
	"context"
	"fmt"

	"github.com/dro14/yordamchi/lib/e"
)

func (c *Cache) isBlocked(ctx context.Context, id string) (bool, error) {

	requests, err := c.Redis.Get(ctx, "blocked:"+id).Int()
	if err != nil {
		if err.Error() == e.KeyNotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	if requests == -14 {
		return true, nil
	} else {
		return false, fmt.Errorf("invalid value: %d", requests)
	}
}

func (c *Cache) isPremium(ctx context.Context, id string) (bool, error) {

	requests, err := c.Redis.Get(ctx, "premium:"+id).Int()
	if err != nil {
		if err.Error() == e.KeyNotFound {
			return false, nil
		} else {
			return false, err
		}
	}

	if requests > 0 {
		return true, nil
	} else {
		return false, fmt.Errorf("invalid value: %d", requests)
	}
}

func (c *Cache) isFree(ctx context.Context, id string) (bool, error) {

	requests, err := c.Redis.Get(ctx, "free:"+id).Int()
	if err != nil {
		if err.Error() == e.KeyNotFound {
			c.Redis.Set(ctx, "free:"+id, NumOfFreeRequests, untilMidnight())
			return true, nil
		} else {
			return false, err
		}
	}

	if requests > 0 && requests <= NumOfFreeRequests {
		return true, nil
	} else if requests == -1 {
		return true, nil
	} else if requests == 0 {
		return false, nil
	} else {
		return false, fmt.Errorf("invalid value: %d", requests)
	}
}
