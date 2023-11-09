package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dro14/yordamchi/lib/types"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IncrementActivity(ctx context.Context, message *tgbotapi.Message, isPremium string) int {
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
	jsonData, err := json.Marshal(activity)
	if err != nil {
		log.Printf("can't encode activity for %s: %v", id(ctx), err)
		return 0
	}
	Client.Set(ctx, "activity:"+id(ctx), string(jsonData), 0)

	keys, err := Client.Keys(ctx, "activity:*").Result()
	if err != nil {
		log.Printf("can't get \"activity:*\": %v", err)
		return 0
	}
	return len(keys)
}

func DecrementActivity(ctx context.Context) {
	Client.Del(ctx, "activity:"+id(ctx))
}

func LoadActivity(ctx context.Context) []*types.Activity {
	var activities []*types.Activity
	keys, err := Client.Keys(ctx, "activity:*").Result()
	if err != nil {
		log.Printf("can't get \"activity:*\": %v", err)
		return activities
	}
	for _, key := range keys {
		result, err := Client.Get(ctx, key).Result()
		if err != nil {
			log.Printf("can't get %q: %v", key, err)
			continue
		}
		activity := &types.Activity{}
		err = json.Unmarshal([]byte(result), activity)
		if err != nil {
			log.Printf("can't decode %q: %v", key, err)
			continue
		}
		activities = append(activities, activity)
	}
	return activities
}
