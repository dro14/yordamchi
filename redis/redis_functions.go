package redis

import (
	"context"
	"fmt"

	"github.com/dro14/yordamchi/lib/e"
)

func isBlocked(ctx context.Context, id string) (bool, error) {

	value, err := Client.Get(ctx, "blocked:"+id).Int()
	if err != nil {
		if err.Error() == e.KeyNotFound {
			return false, nil
		}
		return false, err
	}

	if value == -14 {
		return true, nil
	} else {
		return false, fmt.Errorf("invalid value: %d", value)
	}
}

func isPremium(ctx context.Context, id string) (bool, error) {

	_, err := Client.Get(ctx, "premium:"+id).Result()
	if err != nil {
		if err.Error() == e.KeyNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func isFree(ctx context.Context, id string) (bool, error) {

	requests, err := Client.Get(ctx, "free:"+id).Int()
	if err != nil {
		if err.Error() == e.KeyNotFound {
			err = Client.Set(ctx, "free:"+id, NumOfFreeRequests, untilMidnight()).Err()
			if err != nil {
				return false, err
			}
			return true, nil
		}
		return false, err
	}

	if requests > 0 && requests <= NumOfFreeRequests {
		return true, nil
	} else if requests == 0 {
		return false, nil
	} else {
		return false, fmt.Errorf("invalid number of requests: %d", requests)
	}
}
