package recovery

import (
	"context"
	"log"
	"time"

	"github.com/dro14/yordamchi/client/telegram"
	processor "github.com/dro14/yordamchi/processor/telegram"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init() {

	activities := redis.LoadActivity(context.Background())

	for _, activity := range activities {

		ctx := context.Background()
		ctx = context.WithValue(ctx, "beginning", time.Now())
		ctx = context.WithValue(ctx, "date", activity.Date)
		ctx = context.WithValue(ctx, "user_id", activity.UserID)
		ctx = context.WithValue(ctx, "language_code", activity.LanguageCode)
		ctx = context.WithValue(ctx, "model", redis.Model(ctx))

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

		_, err := telegram.SendMessage(ctx, text.Error[activity.LanguageCode], 0, nil)
		if err != nil {
			log.Printf("can't send error message: %v", err)
		}

		go processor.Process(ctx, message, activity.IsPremium)
	}
}
