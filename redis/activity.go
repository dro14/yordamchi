package redis

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"

	"github.com/dro14/yordamchi/lib/types"
)

func IncrementActivity(ctx context.Context, message *tgbotapi.Message, isPremium string) int {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	activity := &types.Activity{
		MessageID:    message.MessageID,
		Message:      message.Text,
		Date:         message.Date,
		UserID:       message.From.ID,
		FirstName:    ,
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
