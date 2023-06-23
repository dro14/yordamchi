package redis

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/dro14/yordamchi/lib/types"
	"github.com/gotd/td/tg"
)

func IncrementActivity(ctx context.Context, message *tg.Message, user *tg.User, isPremium bool) int {

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

	err = Client.Set(ctx, "activity:"+id, string(data), 0).Err()
	if err != nil {
		log.Printf("can't store activity: %v", err)
	}

	keys, err := Client.Keys(ctx, "activity:*").Result()
	if err != nil {
		log.Printf("can't get keys: %v", err)
		return 0
	}

	return len(keys)
}

func DecrementActivity(ctx context.Context) {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	err := Client.Del(ctx, "activity:"+id).Err()
	if err != nil {
		log.Printf("can't delete activity: %v", err)
	}
}

func LoadActivity(ctx context.Context) []*types.Activity {

	var activities []*types.Activity

	keys, err := Client.Keys(ctx, "activity:*").Result()
	if err != nil {
		log.Printf("can't get keys: %v", err)
		return activities
	}

	for _, key := range keys {

		result, err := Client.Get(ctx, key).Result()
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
