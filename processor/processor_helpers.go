package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dro14/yordamchi/processor/text"
	"github.com/dro14/yordamchi/storage/redis"
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
	time.Sleep(10 * time.Second)
	blockedUsers.Delete(userID)
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

func (p *Processor) messageUpdate(ctx context.Context, message *tgbotapi.Message) (context.Context, bool, bool) {
	if message.From.IsBot || message.Chat.Type != "private" || isBlocked(message.From.ID) {
		return ctx, false, true
	}
	ctx = context.WithValue(ctx, "beginning", time.Now())
	ctx = context.WithValue(ctx, "date", message.Date)
	ctx = context.WithValue(ctx, "user_id", message.From.ID)
	ctx = context.WithValue(ctx, "model", p.redis.Model(ctx))
	ctx = context.WithValue(ctx, "stream", true)
	ctx, foundLang := p.redis.Lang(ctx, message.From.LanguageCode)
	return ctx, true, foundLang
}

func (p *Processor) msg(ctx context.Context) string {
	switch p.redis.UserStatus(ctx) {
	case redis.GPT4Status:
		return fmt.Sprintf(text.Settings2[lang(ctx)], p.redis.GPT4Tokens(ctx))
	case redis.PremiumStatus:
		return fmt.Sprintf(text.Settings1[lang(ctx)], text.PremiumTariff[lang(ctx)], text.Unlimited[lang(ctx)], p.redis.Expiration(ctx))
	default:
		return fmt.Sprintf(text.Settings1[lang(ctx)], text.FreeTariff[lang(ctx)], p.redis.Requests(ctx), p.redis.Expiration(ctx))
	}
}
