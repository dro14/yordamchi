package redis

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"
)

func (r *Redis) SoonExpires(ctx context.Context, pattern string, notifyInterval time.Duration) []int64 {
	var userIDs []int64
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		log.Printf("can't get %q: %s", pattern, err)
		return userIDs
	}

	for _, key := range keys {
		ttl, err := r.client.TTL(ctx, key).Result()
		if err != nil {
			log.Printf("can't get %q: %s", key, err)
			continue
		}
		if 0 < ttl && ttl < notifyInterval {
			_, ID, _ := strings.Cut(key, ":")
			userID, _ := strconv.ParseInt(ID, 10, 64)
			userIDs = append(userIDs, userID)
		}
	}

	return userIDs
}

func (r *Redis) NotifyInterval(ctx context.Context) time.Duration {
	seconds, err := r.client.Get(ctx, "notify_interval").Float64()
	if err != nil {
		log.Printf("can't get %q: %s", "notify_interval", err)
		seconds = 3600.0
	}
	return time.Duration(seconds) * time.Second
}

func (r *Redis) SetNotifyInterval(ctx context.Context, duration time.Duration) {
	r.client.Set(ctx, "notify_interval", duration.Seconds(), 0)
}
