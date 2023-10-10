package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dro14/yordamchi/lib/types"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IncrementActivity(ctx context.Context, message *tgbotapi.Message, isPremium string) int {

	key := fmt.Sprintf("activity:%d", ctx.Value("user_id").(int64))

	activity := &types.Activity{
		MessageID:    message.MessageID,
		Message:      message.Text,
		Date:         message.Date,
		UserID:       message.From.ID,
		FirstName:    message.From.FirstName,
		LastName:     message.From.LastName,
		Username:     message.From.UserName,
		LanguageCode: ctx.Value("language_code").(string),
		IsPremium:    isPremium,
	}

	data, err := json.Marshal(activity)
	if err != nil {
		log.Printf("can't encode activity: %v", err)
		return 0
	}

	err = Client.Set(ctx, key, string(data), 0).Err()
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

	key := fmt.Sprintf("activity:%d", ctx.Value("user_id").(int64))

	err := Client.Del(ctx, key).Err()
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
