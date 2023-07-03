package telegram

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var blockedUsers sync.Map

func isBlocked(userID int64) bool {

	_, ok := blockedUsers.Load(userID)
	if ok {
		return true
	}

	blockedUsers.Store(userID, true)
	go unblockUser(userID)
	return false
}

func unblockUser(userID int64) {
	time.Sleep(1 * time.Minute)
	blockedUsers.Delete(userID)
}

func messageUpdate(ctx context.Context, message *tgbotapi.Message) (context.Context, bool) {

	switch {
	case message.From.IsBot,
		message.Chat.Type != "private",
		isBlocked(message.From.ID):
		return ctx, false
	}

	ctx = context.WithValue(ctx, "beginning", time.Now())
	ctx = context.WithValue(ctx, "date", message.Date)
	ctx = context.WithValue(ctx, "user_id", message.From.ID)
	ctx = context.WithValue(ctx, "language_code", functions.LanguageCode(message.From.LanguageCode))
	ctx = context.WithValue(ctx, "model", redis.Model(ctx))
	return ctx, true
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

func format(timestamp string) string {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return timestamp
	}
	return t.Format(time.DateTime)
}

func slice(completion string) []string {

	var completions []string

	for len(completion) > 4096 {
		cutIndex := 0
	Loop:
		for i := 4096; i >= 0; i-- {
			switch completion[i] {
			case ' ', '\n', '\t', '\r':
				cutIndex = i
				break Loop
			}
		}
		completions = append(completions, completion[:cutIndex])
		completion = completion[cutIndex:]
	}

	return append(completions, completion)
}

func msg(ctx context.Context, lang string) string {

	switch redis.UserStatus(ctx) {
	case types.GPT4Status:
		return fmt.Sprintf(text.Settings2[lang], redis.GPT4Tokens(ctx))
	case types.PremiumStatus:
		return fmt.Sprintf(text.Settings1[lang], text.PremiumTariff[lang], text.Unlimited[lang], redis.Expiration(ctx))
	default:
		return fmt.Sprintf(text.Settings1[lang], text.FreeTariff[lang], text.Unlimited[lang], redis.Expiration(ctx))
	}
}
