package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/go-redis/redis/v8"
	"github.com/gotd/td/tg"
	"log"
	"os"
	"strconv"
)

const NumOfFreeRequests = 10

type Cache struct {
	Redis *redis.Client
}

func New() *Cache {

	url, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		log.Fatalf("redis url is not specified")
	}

	password, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		log.Fatalf("redis password is not specified")
	}

	return &Cache{
		Redis: redis.NewClient(
			&redis.Options{
				Addr:     url,
				Password: password,
				DB:       0,
			},
		),
	}
}

func (c *Cache) Status(ctx context.Context) (types.UserStatus, error) {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	result, err := c.isBlocked(ctx, id)
	if err != nil {
		log.Printf("can't check whether user %s is blocked: %v", id, err)
		return types.UnknownStatus, err
	} else if result {
		return types.BlockedStatus, nil
	}

	result, err = c.isPremium(ctx, id)
	if err != nil {
		log.Printf("can't check whether user %s is premium: %v", id, err)
		return types.UnknownStatus, err
	} else if result {
		return types.PremiumStatus, nil
	}

	result, err = c.isFree(ctx, id)
	if err != nil {
		log.Printf("can't check whether user %s is free: %v", id, err)
		return types.UnknownStatus, err
	} else if result {
		return types.FreeStatus, nil
	}

	return types.ExhaustedStatus, nil
}

func (c *Cache) Balance(ctx context.Context) (int, error) {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	requests, err := c.Redis.Get(ctx, "premium:"+id).Int()
	if err == nil {
		return requests, nil
	} else if err.Error() != e.KeyNotFound {
		log.Printf("can't get \"premium:%s\": %v", id, err)
		return -1, err
	}

	requests, err = c.Redis.Get(ctx, "free:"+id).Int()
	if err == nil {
		return requests, nil
	} else if err.Error() != e.KeyNotFound {
		log.Printf("can't get \"free:%s\": %v", id, err)
		return -1, err
	}

	log.Printf(e.UserNotDefined)
	return -1, e.UserNotDefinedError
}

func (c *Cache) Decrement(ctx context.Context) error {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	requests, err := c.Redis.Get(ctx, "premium:"+id).Int()
	if err == nil {
		if requests > 0 {
			if requests == 1 {
				c.Redis.Del(ctx, "premium:"+id)
			} else {
				c.Redis.Set(ctx, "premium:"+id, requests-1, 0)
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

	requests, err = c.Redis.Get(ctx, "free:"+id).Int()
	if err == nil {
		if requests > 0 && requests <= NumOfFreeRequests {
			c.Redis.Set(ctx, "free:"+id, requests-1, untilMidnight())
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

func (c *Cache) LoadContext(ctx context.Context, prompt string) []types.Message {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	messages := make([]types.Message, 0, 3)
	message := types.Message{
		Content: prompt,
		Role:    "user",
	}

	result, err := c.Redis.Get(ctx, "context:"+id).Result()
	if err != nil {
		if err.Error() != e.KeyNotFound {
			log.Printf("can't get \"context:%s\": %v", id, err)
		}
		return append(messages, message)
	}

	err = json.Unmarshal([]byte(result), &messages)
	if err != nil {
		log.Printf("can't decode context: %v", err)
		return append(messages, message)
	}

	return append(messages, message)
}

func (c *Cache) StoreContext(ctx context.Context, prompt, completion string) {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	messages := []types.Message{
		{
			Content: prompt,
			Role:    "user",
		},
		{
			Content: completion,
			Role:    "assistant",
		},
	}

	data, err := json.Marshal(messages)
	if err != nil {
		log.Printf("can't encode context: %v", err)
		return
	}

	err = c.Redis.Set(ctx, "context:"+id, string(data), 0).Err()
	if err != nil {
		log.Printf("can't store context: %v", err)
	}
}

func (c *Cache) DeleteContext(ctx context.Context) {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	err := c.Redis.Del(ctx, "context:"+id).Err()
	if err != nil {
		log.Printf("can't delete context: %v", err)
	}
}

func (c *Cache) IncrementActivity(ctx context.Context, message *tg.Message, user *tg.User, isPremium bool) int {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	activity := &types.Activity{
		MessageID:    message.ID,
		Message:      message.Message,
		Date:         message.Date,
		UserID:       user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		LanguageCode: user.LangCode,
		IsPremium:    isPremium,
	}

	data, err := json.Marshal(activity)
	if err != nil {
		log.Printf("can't encode activity: %v", err)
		return 0
	}

	err = c.Redis.Set(ctx, "activity:"+id, string(data), 0).Err()
	if err != nil {
		log.Printf("can't store activity: %v", err)
	}

	keys, err := c.Redis.Keys(ctx, "activity:*").Result()
	if err != nil {
		log.Printf("can't get keys: %v", err)
		return 0
	}

	return len(keys)
}

func (c *Cache) DecrementActivity(ctx context.Context) {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	err := c.Redis.Del(ctx, "activity:"+id).Err()
	if err != nil {
		log.Printf("can't delete activity: %v", err)
	}
}

func (c *Cache) LoadActivity(ctx context.Context) []*types.Activity {

	var activities []*types.Activity

	keys, err := c.Redis.Keys(ctx, "activity:*").Result()
	if err != nil {
		log.Printf("can't get keys: %v", err)
		return activities
	}

	for _, key := range keys {

		result, err := c.Redis.Get(ctx, key).Result()
		if err != nil {
			log.Printf("can't get \"%s\": %v", key, err)
			continue
		}

		activity := &types.Activity{}
		err = json.Unmarshal([]byte(result), activity)
		if err != nil {
			log.Printf("can't decode activity: %v", err)
			continue
		}

		activities = append(activities, activity)
	}

	return activities
}
