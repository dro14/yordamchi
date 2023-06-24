package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dro14/yordamchi/lib/types"
	"github.com/go-redis/redis/v8"
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

func UserStatus(ctx context.Context) types.UserStatus {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	result, err := isBlocked(ctx, id)
	if err != nil {
		log.Printf("can't check whether user %s is blocked: %v", id, err)
		return types.UnknownStatus
	} else if result {
		return types.BlockedStatus
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

func Decrement(ctx context.Context) {

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

func SetPremium(userID int64, amount int, Type string) error {

	key := fmt.Sprintf("premium:%d", userID)

	var expiration time.Time
	switch Type {
	case "weekly", "requests":
		expiration = time.Now().AddDate(0, 0, 7)
	case "monthly":
		expiration = time.Now().AddDate(0, 1, 0)
	default:
		return fmt.Errorf("invalid type: %s", Type)
	}

	err := Client.Set(context.Background(), key, expiration.Format("15:04:05 02.01.2006"), time.Until(expiration)).Err()
	if err != nil {
		log.Printf("can't set premium: %v", err)
		return err
	}

	return nil
}
