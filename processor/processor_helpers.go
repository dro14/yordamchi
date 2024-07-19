package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/dro14/yordamchi/storage/redis/status"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var messageBuffer sync.Map

func start(ctx context.Context) time.Time {
	return ctx.Value("start").(time.Time)
}

func userStatus(ctx context.Context) status.Status {
	return ctx.Value("user_status").(status.Status)
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

func (p *Processor) messageUpdate(ctx context.Context, message *tgbotapi.Message) (context.Context, bool, bool) {
	value, ok := messageBuffer.Load(message.From.ID)
	if ok {
		message.Text = value.(string) + "\n" + message.Text
	}
	messageBuffer.Store(message.From.ID, message.Text)
	time.Sleep(200 * time.Millisecond)
	value, _ = messageBuffer.Load(message.From.ID)

	if message.From.IsBot || message.Chat.Type != "private" || value != message.Text {
		return ctx, true, true
	}
	ctx = context.WithValue(ctx, "start", time.Now())
	ctx = context.WithValue(ctx, "user_id", message.From.ID)
	ctx = context.WithValue(ctx, "stream", true)
	ctx = context.WithValue(ctx, "json_mode", false)
	ctx = context.WithValue(ctx, "user_status", p.redis.UserStatus(ctx))
	if userStatus(ctx) == status.Premium {
		ctx = context.WithValue(ctx, "model", models.GPT4o)
	} else {
		ctx = context.WithValue(ctx, "model", models.GPT4oMini)
	}
	ctx, foundLang := p.redis.Lang(ctx, message.From.LanguageCode)
	return ctx, false, foundLang
}

func (p *Processor) msg(ctx context.Context) string {
	switch userStatus(ctx) {
	case status.Premium:
		template := text.Settings2[lang(ctx)]
		requests, expiration := p.redis.Requests(ctx), p.redis.Expiration(ctx)
		return fmt.Sprintf(template, requests, expiration)
	case status.Unlimited:
		template := text.Settings1[lang(ctx)]
		expiration := p.redis.Expiration(ctx)
		return fmt.Sprintf(template, expiration)
	default:
		template := text.Settings[lang(ctx)]
		requests, expiration := p.redis.Requests(ctx), p.redis.Expiration(ctx)
		return fmt.Sprintf(template, requests, expiration)
	}
}
