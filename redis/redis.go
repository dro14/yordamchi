package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dro14/yordamchi/lib/types"
	"github.com/go-redis/redis/v8"
)

const (
	NumOfFreeRequests = 5
	KeyNotFound       = "redis: nil"
)

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

func UserStatus(ctx context.Context) types.UserStatus {

	id := fmt.Sprintf("%d", ctx.Value("user_id").(int64))

	result, err := isGPT4(ctx, id)
	if err != nil {
		log.Printf("can't check whether user %s is gpt-4: %v", id, err)
		return types.UnknownStatus
	} else if result {
		return types.GPT4Status
	}

	result, err = isPremium(ctx, id)
	if err != nil {
		log.Printf("can't check whether user %s is premium: %v", id, err)
		return types.UnknownStatus
	} else if result {
		return types.PremiumStatus
	}

	result, err = isFree(ctx, id)
	if err != nil {
		log.Printf("can't check whether user %s is free: %v", id, err)
		return types.UnknownStatus
	} else if result {
		return types.FreeStatus
	}

	return types.ExhaustedStatus
}

func Expiration(ctx context.Context) string {

	key := fmt.Sprintf("premium:%d", ctx.Value("user_id").(int64))

	value, err := Client.Get(ctx, key).Result()
	if err != nil {
		return midnight()
	}

	return value
}

func Requests(ctx context.Context) string {

	key := fmt.Sprintf("free:%d", ctx.Value("user_id").(int64))

	requests, err := Client.Get(ctx, key).Int()
	if err != nil {
		log.Printf("can't get %q: %v", key, err)
		return ""
	}

	return fmt.Sprintf("%d/%d", requests, NumOfFreeRequests)
}

func Decrement(ctx context.Context, used int) {

	if ctx.Value("model") == "gpt-4" {
		key := fmt.Sprintf("gpt-4:%d", ctx.Value("user_id").(int64))

		available, err := Client.Get(ctx, key).Int()
		if err != nil {
			log.Printf("can't get %q: %v", key, err)
			return
		}

		if available <= used {
			err = Client.Del(ctx, key).Err()
			if err != nil {
				log.Printf("can't delete %q: %v", key, err)
			}
		} else {
			err = Client.Set(ctx, key, available-used, 0).Err()
			if err != nil {
				log.Printf("can't decrement %q: %v", key, err)
			}
		}
	} else {
		key := fmt.Sprintf("premium:%d", ctx.Value("user_id").(int64))

		_, err := Client.Get(ctx, key).Result()
		if err == nil {
			return
		}

		key = fmt.Sprintf("free:%d", ctx.Value("user_id").(int64))

		requests, err := Client.Get(ctx, key).Int()
		if err != nil {
			log.Printf("can't get %q: %v", key, err)
			return
		}

		if requests > 0 && requests <= NumOfFreeRequests {
			err = Client.Set(ctx, key, requests-1, untilMidnight()).Err()
			if err != nil {
				log.Printf("can't decrement %q: %v", key, err)
			}
		} else {
			log.Printf("invalid number of requests: %d", requests)
		}
	}
}

func PerformTransaction(userID int64, amount int, Type string) {

	ctx := context.Background()

	if Type == "gpt-4" {
		key := fmt.Sprintf("gpt-4:%d", userID)

		tokens, _ := Client.Get(ctx, key).Int()
		tokens += amount / 100

		Client.Set(ctx, key, tokens, 0)
	} else {
		key := fmt.Sprintf("premium:%d", userID)

		var expiration time.Time
		switch Type {
		case "weekly":
			expiration = time.Now().AddDate(0, 0, 7)
		case "monthly":
			expiration = time.Now().AddDate(0, 1, 0)
		}

		expirationDate := expiration.Format("15:04:05 02.01.2006")
		Client.Set(ctx, key, expirationDate, time.Until(expiration))
	}
}
