package telegram

import (
	"context"
	"log"
	"time"

	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
	"github.com/gotd/td/tg"
)

func (p *Processor) Recover() {

	activities := redis.LoadActivity(context.Background())

	for _, activity := range activities {

		start := time.Now()
		ctx := context.Background()
		ctx = context.WithValue(ctx, "start", start)
		ctx = context.WithValue(ctx, "date", activity.Date)
		ctx = context.WithValue(ctx, "user_id", activity.UserID)
		ctx = context.WithValue(ctx, "language_code", functions.LanguageCode(activity.LanguageCode))
		ctx = context.WithValue(ctx, "model", redis.Model(ctx))

		message := &tg.Message{
			ID:      activity.MessageID,
			Message: activity.Message,
		}

		user := &tg.User{
			ID:        activity.UserID,
			FirstName: activity.FirstName,
			LastName:  activity.LastName,
			Username:  activity.Username,
			LangCode:  activity.LanguageCode,
		}

		_, err := p.Client.SendMessage(ctx, text.Error[lang(ctx)], 0, nil)
		if err != nil {
			log.Printf("can't send error message: %v", err)
		}

		go p.Stream(ctx, message, user, activity.IsPremium)
	}
}
