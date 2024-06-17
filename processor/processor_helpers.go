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

func model(ctx context.Context) string {
	return ctx.Value("model").(string)
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
	ctx = context.WithValue(ctx, "translate", false)
	ctx = context.WithValue(ctx, "json_mode", false)
	ctx = context.WithValue(ctx, "user_status", p.redis.UserStatus(ctx))
	if userStatus(ctx) == status.Premium {
		ctx = context.WithValue(ctx, "model", models.GPT4)
	} else {
		ctx = context.WithValue(ctx, "model", models.GPT3)
	}
	ctx, foundLang := p.redis.Lang(ctx, message.From.LanguageCode)
	return ctx, false, foundLang
}

func (p *Processor) msg(ctx context.Context) string {
	switch userStatus(ctx) {
	case status.Premium:
		template := text.Settings2[lang(ctx)]
		expiration, requests := p.redis.Premium(ctx)
		return fmt.Sprintf(template, requests, expiration)
	case status.Unlimited:
		template := text.Settings1[lang(ctx)]
		return fmt.Sprintf(template, p.redis.Expiration(ctx))
	default:
		template := text.Settings[lang(ctx)]
		return fmt.Sprintf(template, p.redis.Requests(ctx))
	}
}

func (p *Processor) needTranslation(ctx context.Context) bool {
	return model(ctx) == models.GPT3 && lang(ctx) == "uz"
}
