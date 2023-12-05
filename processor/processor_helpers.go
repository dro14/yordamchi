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

func start(ctx context.Context) time.Time {
	return ctx.Value("start").(time.Time)
}

func userStatus(ctx context.Context) redis.UserStatus {
	return ctx.Value("user_status").(redis.UserStatus)
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

func model(ctx context.Context) string {
	return ctx.Value("model").(string)
}

func (p *Processor) messageUpdate(ctx context.Context, message *tgbotapi.Message) (context.Context, bool, bool) {
	if message.From.IsBot || message.Chat.Type != "private" || isBlocked(message.From.ID) {
		return ctx, true, true
	}
	ctx = context.WithValue(ctx, "start", time.Now())
	ctx = context.WithValue(ctx, "user_id", message.From.ID)
	ctx = context.WithValue(ctx, "user_status", p.redis.UserStatus(ctx))
	ctx = context.WithValue(ctx, "stream", true)
	ctx, foundLang := p.redis.Lang(ctx, message.From.LanguageCode)
	return ctx, false, foundLang
}

func (p *Processor) msg(ctx context.Context) string {
	switch userStatus(ctx) {
	case redis.StatusPremium:
		template := text.Settings2[lang(ctx)]
		return fmt.Sprintf(template, p.redis.Expiration(ctx))
	case redis.StatusUnlimited:
		template := text.Settings1[lang(ctx)]
		return fmt.Sprintf(template, p.redis.Expiration(ctx))
	default:
		template := text.Settings[lang(ctx)]
		return fmt.Sprintf(template, p.redis.Requests(ctx), p.redis.Expiration(ctx))
	}
}
