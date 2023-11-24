package processor

import (
	"context"
	"log"

	"github.com/dro14/yordamchi/clients/openai"
	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/clients/other"
	"github.com/dro14/yordamchi/clients/telegram"
	"github.com/dro14/yordamchi/payme"
	"github.com/dro14/yordamchi/storage/postgres"
	"github.com/dro14/yordamchi/storage/redis"
	"github.com/dro14/yordamchi/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Processor struct {
	telegram *telegram.Telegram
	postgres *postgres.Postgres
	openai   *openai.OpenAI
	redis    *redis.Redis
	payme    *payme.Payme
	apis     *other.APIs
}

func New() *Processor {
	return &Processor{
		telegram: telegram.New(),
		postgres: postgres.New(),
		openai:   openai.New(),
		redis:    redis.New(),
		payme:    payme.New(),
		apis:     other.New(),
	}
}

func (p *Processor) Update(update *tgbotapi.Update) {
	defer utils.RecoverIfPanic()
	ctx := context.Background()
	switch {
	case update.Message != nil:
		p.Message(ctx, update.Message)
	case update.CallbackQuery != nil:
		p.CallbackQuery(ctx, update.CallbackQuery)
	case update.MyChatMember != nil:
		p.MyChatMember(ctx, update.MyChatMember)
	default:
		log.Printf("unknown update type:\n%+v", update)
	}
}

func (p *Processor) Message(ctx context.Context, message *tgbotapi.Message) {
	ctx, allowed, foundLang := p.messageUpdate(ctx, message)
	if !allowed || !foundLang {
		if !foundLang {
			p.language(ctx)
			blockedUsers.Delete(message.From.ID)
		}
		return
	}
	defer blockedUsers.Delete(message.From.ID)

	done := p.doCommand(ctx, message)
	if done {
		return
	}

	switch ctx.Value("user_status") {
	case redis.StatusPremium:
		ctx = context.WithValue(ctx, "model", models.GPT4)
		if message.Text != "" || message.Photo != nil {
			p.Process(ctx, message, "true")
		}
	case redis.StatusFree:
		ctx = context.WithValue(ctx, "model", models.GPT3)
		if message.Text != "" {
			p.Process(ctx, message, "false")
		} else if message.Photo != nil {
			p.premiumFeature(ctx)
		}
	case redis.StatusExhausted:
		p.exhausted(ctx)
	default:
		log.Println("unknown user status:", message.From.ID)
	}
}

func (p *Processor) CallbackQuery(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	ctx = context.WithValue(ctx, "date", callbackQuery.Message.Date)
	ctx = context.WithValue(ctx, "user_id", callbackQuery.From.ID)
	ctx, _ = p.redis.Lang(ctx, callbackQuery.From.LanguageCode)

	switch callbackQuery.Data {
	case "new_chat":
		p.newChatCallback(ctx, callbackQuery)
	case "help":
		p.helpCallback(ctx, callbackQuery)
	case "settings":
		p.settingsCallback(ctx, callbackQuery)
	case "uz", "ru", "en":
		p.languageCallback(ctx, callbackQuery)
	case "examples":
		p.examplesCallback(ctx, callbackQuery)
	case "vivid", "natural":
		p.generateCallback(ctx, callbackQuery)
	default:
		log.Println("unknown callback data:", callbackQuery.Data)
	}
}

func (p *Processor) MyChatMember(ctx context.Context, chatMemberUpdated *tgbotapi.ChatMemberUpdated) {
	ctx = context.WithValue(ctx, "date", chatMemberUpdated.Date)
	ctx = context.WithValue(ctx, "user_id", chatMemberUpdated.From.ID)
	ctx, _ = p.redis.Lang(ctx, chatMemberUpdated.From.LanguageCode)

	switch chatMemberUpdated.NewChatMember.Status {
	case "kicked":
		p.postgres.UserBlocked(ctx, &chatMemberUpdated.From)
	case "member":
		p.postgres.UserRestarted(ctx, &chatMemberUpdated.From)
	default:
		log.Println("unknown chat member status:", chatMemberUpdated.NewChatMember.Status)
	}
}
