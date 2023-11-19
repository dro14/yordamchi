package recover

import (
	"context"
	"log"
	"time"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/clients/telegram"
	"github.com/dro14/yordamchi/processor"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/dro14/yordamchi/storage/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start() {
	tgProcessor := processor.New()
	tgClient := telegram.New()
	redisClient := redis.New()

	for _, activity := range redisClient.Activities() {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "start", time.Now())
		ctx = context.WithValue(ctx, "user_id", activity.UserID)
		ctx = context.WithValue(ctx, "user_status", redisClient.UserStatus(ctx))
		ctx = context.WithValue(ctx, "stream", true)
		ctx = context.WithValue(ctx, "language_code", activity.LanguageCode)

		if ctx.Value("user_status") == redis.StatusPremium {
			ctx = context.WithValue(ctx, "model", models.GPT4)
		} else {
			ctx = context.WithValue(ctx, "model", models.GPT3)
		}

		message := &tgbotapi.Message{
			MessageID: activity.MessageID,
			Text:      activity.Message,
			Date:      activity.Date,
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
