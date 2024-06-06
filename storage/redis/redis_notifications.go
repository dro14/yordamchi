package redis

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dro14/yordamchi/utils"
)

func (r *Redis) SoonExpires(ctx context.Context, pattern string) []int64 {
	var userIDs []int64
	keys, err := client.Keys(ctx, pattern).Result()
	if err != nil {
		log.Printf("can't get %q: %s", pattern, err)
		return userIDs
	}

	for _, key := range keys {
		ttl, err := client.TTL(ctx, key).Result()
		if err != nil {
			log.Printf("can't get %q: %s", key, err)
			continue
		}
		if 0 < ttl && ttl < utils.NotificationInterval {
			_, ID, _ := strings.Cut(key, ":")
			userID, _ := strconv.ParseInt(ID, 10, 64)
			userIDs = append(userIDs, userID)
		}
	}

	if pattern == "lang:*" || pattern == "context:*" {
		for _, userID := range userIDs {
			client.IncrBy(ctx, fmt.Sprintf("free:%d", userID), 5)
		}
	}

	return userIDs
}

func (r *Redis) NotificationTime(ctx context.Context) time.Time {
	unixTime, err := client.Get(ctx, "notification_time").Int64()
	if err != nil {
		log.Printf("can't get %q: %s", "notification_time", err)
		return time.Now().Add(utils.NotificationInterval)
	}
	return time.Unix(unixTime, 0)
}

func (r *Redis) SetNotificationTime(ctx context.Context, notificationTime time.Time) {
	client.Set(ctx, "notification_time", notificationTime.Unix(), 0)
}
