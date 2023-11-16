package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (r *Redis) UserActivity() []*Activity {
	ctx := context.Background()
	var activities []*Activity
	keys, err := r.client.Keys(ctx, "activity:*").Result()
	if err != nil {
		log.Println("can't get \"activity:*\":", err)
		return activities
	}
	var result []byte
	for _, key := range keys {
		result, err = r.client.Get(ctx, key).Bytes()
		if err != nil {
			log.Printf("can't get %q: %s", key, err)
			continue
		}
		activity := &Activity{}
		err = json.Unmarshal(result, activity)
		if err != nil {
			log.Printf("can't decode %q: %s", key, err)
			continue
		}
		activities = append(activities, activity)
	}
	return activities
}

func (r *Redis) IncrementActivity(ctx context.Context, message *tgbotapi.Message, isPremium string) int {
	activity := &Activity{
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
		log.Printf("user %s: can't encode activity: %s", id(ctx), err)
		return 0
	}
	r.client.Set(ctx, "activity:"+id(ctx), jsonData, 0)

	keys, err := r.client.Keys(ctx, "activity:*").Result()
	if err != nil {
		log.Println("can't get \"activity:*\":", err)
		return 0
	}
	return len(keys)
}

func (r *Redis) DecrementActivity(ctx context.Context) {
	r.client.Del(ctx, "activity:"+id(ctx))
}
