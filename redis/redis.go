package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/models"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/go-redis/redis/v8"
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
	options := &redis.Options{Addr: url, Password: password}
	Client = redis.NewClient(options)
}

func UserStatus(ctx context.Context) types.UserStatus {
	result, err := isGPT4(ctx)
	if err != nil {
		log.Printf("can't check whether user %s status is gpt-4: %v", id(ctx), err)
		return types.UnknownStatus
	} else if result {
		return types.GPT4Status
	}

	result, err = isPremium(ctx)
	if err != nil {
		log.Printf("can't check whether user %s status is premium: %v", id(ctx), err)
		return types.UnknownStatus
	} else if result {
		return types.PremiumStatus
	}

	result, err = isFree(ctx)
	if err != nil {
		log.Printf("can't check whether user %s status is free: %v", id(ctx), err)
		return types.UnknownStatus
	} else if result {
		return types.FreeStatus
	}

	return types.ExhaustedStatus
}

func Expiration(ctx context.Context) string {
	value, err := Client.Get(ctx, "premium:"+id(ctx)).Result()
	if err != nil {
		return midnight()
	}
	return value
}

func Requests(ctx context.Context) string {
	requests, err := Client.Get(ctx, "free:"+id(ctx)).Int()
	if err != nil {
		log.Printf("can't get %q: %v", "free:"+id(ctx), err)
		return ""
	}
	return fmt.Sprintf("%d/%d", requests, constants.NumOfFreeRequests)
}

func Decrement(ctx context.Context, used int) {
	switch ctx.Value("model") {
	case models.GPT4, models.GPT4V:
		available, err := Client.Get(ctx, "gpt-4:"+id(ctx)).Int()
		if err != nil {
			log.Printf("can't get %q: %v", "gpt-4:"+id(ctx), err)
			return
		}
		if available <= used {
			Client.Del(ctx, "gpt-4:"+id(ctx))
		} else {
			Client.Set(ctx, "gpt-4:"+id(ctx), available-used, 0)
		}
	default:
		_, err := Client.Get(ctx, "premium:"+id(ctx)).Result()
		if err == nil {
			return
		}
		requests, err := Client.Get(ctx, "free:"+id(ctx)).Int()
		if err != nil {
			log.Printf("can't get %q: %v", "free:"+id(ctx), err)
			return
		}
		if requests > 0 && requests <= constants.NumOfFreeRequests {
			Client.Set(ctx, "free:"+id(ctx), requests-1, untilMidnight())
		} else {
			log.Printf("invalid number of requests for %s: %d", id(ctx), requests)
		}
	}
}

func PerformTransaction(userID int64, amount int, Type string) {
	ctx := context.Background()
	ID := fmt.Sprintf("%d", userID)
	if Type == "gpt-4" {
		tokens, _ := Client.Get(ctx, "gpt-4:"+ID).Int()
		tokens += amount / 100
		Client.Set(ctx, "gpt-4:"+ID, tokens, 0)
	} else {
		var expiration time.Time
		switch Type {
		case "weekly":
			expiration = time.Now().AddDate(0, 0, 7)
		case "monthly":
			expiration = time.Now().AddDate(0, 1, 0)
		}
		expirationDate := expiration.Format("15:04:05 02.01.2006")
		Client.Set(ctx, "premium:"+ID, expirationDate, time.Until(expiration))
	}
}
