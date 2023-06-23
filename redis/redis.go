package redis

import (
	"context"
	"fmt"
	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"strconv"
)

const NumOfFreeRequests = 10

var Client *redis.Client

func Init() {

	url, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		log.Fatalf("redis url is not specified")
	}

	password, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		log.Fatalf("redis password is not specified")
	}

	Client = redis.NewClient(
		&redis.Options{
			Addr:     url,
			Password: password,
			DB:       0,
		},
	)
}

func Status(ctx context.Context) (types.UserStatus, error) {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	result, err := isBlocked(ctx, id)
	if err != nil {
		log.Printf("can't check whether user %s is blocked: %v", id, err)
		return types.UnknownStatus, err
	} else if result {
		return types.BlockedStatus, nil
	}

	result, err = isPremium(ctx, id)
	if err != nil {
		log.Printf("can't check whether user %s is premium: %v", id, err)
		return types.UnknownStatus, err
	} else if result {
		return types.PremiumStatus, nil
	}

	result, err = isFree(ctx, id)
	if err != nil {
		log.Printf("can't check whether user %s is free: %v", id, err)
		return types.UnknownStatus, err
	} else if result {
		return types.FreeStatus, nil
	}

	return types.ExhaustedStatus, nil
}

func Balance(ctx context.Context) (int, error) {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	requests, err := Client.Get(ctx, "premium:"+id).Int()
	if err == nil {
		return requests, nil
	} else if err.Error() != e.KeyNotFound {
		log.Printf("can't get \"premium:%s\": %v", id, err)
		return -1, err
	}

	requests, err = Client.Get(ctx, "free:"+id).Int()
	if err == nil {
		return requests, nil
	} else if err.Error() != e.KeyNotFound {
		log.Printf("can't get \"free:%s\": %v", id, err)
		return -1, err
	}

	log.Printf(e.UserNotDefined)
	return -1, e.UserNotDefinedError
}

func Decrement(ctx context.Context) error {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	requests, err := Client.Get(ctx, "premium:"+id).Int()
	if err == nil {
		if requests > 0 {
			if requests == 1 {
				Client.Del(ctx, "premium:"+id)
			} else {
				Client.Set(ctx, "premium:"+id, requests-1, 0)
			}
			return nil
		} else {
			log.Printf("invalid value: %d", requests)
			return fmt.Errorf("invalid value: %d", requests)
		}
	} else if err.Error() != e.KeyNotFound {
		log.Printf("can't get \"premium:%s\": %v", id, err)
		return err
	}

	requests, err = Client.Get(ctx, "free:"+id).Int()
	if err == nil {
		if requests > 0 && requests <= NumOfFreeRequests {
			Client.Set(ctx, "free:"+id, requests-1, untilMidnight())
			return nil
		} else if requests == -1 {
			return nil
		} else {
			log.Printf("invalid value: %d", requests)
			return fmt.Errorf("invalid value: %d", requests)
		}
	} else if err.Error() != e.KeyNotFound {
		log.Printf("can't get \"free:%s\": %v", id, err)
		return err
	}

	log.Printf(e.UserNotDefined)
	return e.UserNotDefinedError
}

func SetPremium(userID int64) error {

	key := fmt.Sprintf("premium:%d", userID)

	requests, err := Client.Get(context.Background(), key).Int()
	if err == nil {
		err = Client.Set(context.Background(), key, requests+200, 0).Err()
		if err != nil {
			log.Printf("can't set premium: %v", err)
			return err
		}
	} else if err.Error() == e.KeyNotFound {
		err = Client.Set(context.Background(), key, 200, 0).Err()
		if err != nil {
			log.Printf("can't set premium: %v", err)
			return err
		}
	} else {
		log.Printf("can't get premium: %v", err)
		return err
	}

	return nil
}
