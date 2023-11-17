package recover

import (
	"context"
	"log"
	"time"

	"github.com/dro14/yordamchi/clients/telegram"
	"github.com/dro14/yordamchi/processor"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/dro14/yordamchi/storage/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init() {
	tgProcessor := processor.New()
	tgClient := telegram.New()
	redisClient := redis.New()

	activities := redisClient.Activities()
	for _, activity := range activities {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "beginning", time.Now())
		ctx = context.WithValue(ctx, "date", activity.Date)
		ctx = context.WithValue(ctx, "user_id", activity.UserID)
		ctx = context.WithValue(ctx, "language_code", activity.LanguageCode)
		ctx = context.WithValue(ctx, "model", redisClient.Model(ctx))
		ctx = context.WithValue(ctx, "stream", true)

		message := &tgbotapi.Message{
			MessageID: activity.MessageID,
			Text:      activity.Message,
			From: &tgbotapi.User{
				ID:           activity.UserID,
				FirstName:    activity.FirstName,
				LastName:     activity.LastName,
				UserName:     activity.Username,
				LanguageCode: activity.LanguageCode,
			},
		}

		_, err := tgClient.SendMessage(ctx, text.Error[activity.LanguageCode], 0, nil)
		if err != nil {
			log.Println("can't send error message:", err)
			continue
		}
		go tgProcessor.Process(ctx, message, activity.IsPremium)
	}
}
