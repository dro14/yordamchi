package telegram

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/processor/telegram/info_bot"
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

func messageUpdate(ctx context.Context, message *tgbotapi.Message) (context.Context, bool, bool) {
	switch {
	case message.From.IsBot, message.Chat.Type != "private", isBlocked(message.From.ID):
		return ctx, true, false
	}
	ctx = context.WithValue(ctx, "beginning", time.Now())
	ctx = context.WithValue(ctx, "date", message.Date)
	ctx = context.WithValue(ctx, "user_id", message.From.ID)
	ctx = context.WithValue(ctx, "model", redis.Model(ctx))
	ctx, shouldSetLang := redis.Lang(ctx)
	return ctx, false, shouldSetLang
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

func msg(ctx context.Context, lang string) string {
	switch redis.UserStatus(ctx) {
	case types.GPT4Status:
		return fmt.Sprintf(text.Settings2[lang], redis.GPT4Tokens(ctx))
	case types.PremiumStatus:
		return fmt.Sprintf(text.Settings1[lang], text.PremiumTariff[lang], text.Unlimited[lang], redis.Expiration(ctx))
	default:
		return fmt.Sprintf(text.Settings1[lang], text.FreeTariff[lang], redis.Requests(ctx), redis.Expiration(ctx))
	}
}

func recoverFromPanic() {
	if r := recover(); r != nil {
		info_bot.Send(fmt.Sprintf("panic:\n%v", r))
		log.Fatalf("fatal error: restarting bot")
	}
}
