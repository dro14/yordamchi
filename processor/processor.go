package processor

import (
	"context"
	"log"

	"github.com/dro14/yordamchi/clients/openai"
	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/clients/other"
	"github.com/dro14/yordamchi/clients/telegram"
	"github.com/dro14/yordamchi/configs"
	"github.com/dro14/yordamchi/payme"
	"github.com/dro14/yordamchi/storage/postgres"
	"github.com/dro14/yordamchi/storage/redis"
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

func (p *Processor) Message(message *tgbotapi.Message) {
	defer configs.RecoverIfPanic()
	ctx, allowed, foundLang := p.messageUpdate(message)
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

	switch p.redis.UserStatus(ctx) {
	case redis.GPT4Status:
		if p.redis.GPT4Tokens(ctx) > 0 {
			p.Process(ctx, message, models.GPT4)
		} else {
			p.gpt4(ctx)
		}
	case redis.PremiumStatus:
		p.Process(ctx, message, "true")
	case redis.FreeStatus:
		if message.Photo == nil {
			p.Process(ctx, message, "false")
		} else {
			// TODO
		}
	case redis.ExhaustedStatus:
		p.exhausted(ctx)
	default:
		log.Println("unknown user status:", message.From.ID)
	}
}

func (p *Processor) CallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {
	defer configs.RecoverIfPanic()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "date", callbackQuery.Message.Date)
	ctx = context.WithValue(ctx, "user_id", callbackQuery.From.ID)
	ctx, _ = p.redis.Lang(ctx, callbackQuery.From.LanguageCode)

	switch callbackQuery.Data {
	case "new_chat":
		p.newChat(ctx)
	case "examples":
		p.examplesCallback(ctx, callbackQuery.Message.MessageID)
	case "help":
		p.helpCallback(ctx, callbackQuery.Message.MessageID)
	case models.GPT3, models.GPT4:
		p.model(ctx, callbackQuery.Message.MessageID, callbackQuery.Data)
	case "uz", "ru", "en":
		p.languageCallback(ctx, callbackQuery.Message, callbackQuery.Data)
	default:
		log.Println("unknown callback data:", callbackQuery.Data)
	}
}

func (p *Processor) MyChatMember(chatMemberUpdated *tgbotapi.ChatMemberUpdated) {
	defer configs.RecoverIfPanic()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "date", chatMemberUpdated.Date)
	ctx = context.WithValue(ctx, "user_id", chatMemberUpdated.From.ID)
	ctx, _ = p.redis.Lang(ctx, chatMemberUpdated.From.LanguageCode)

	switch chatMemberUpdated.NewChatMember.Status {
	case "kicked":
		p.postgres.DeactivateUser(ctx, &chatMemberUpdated.From)
	case "member":
		p.postgres.RejoinUser(ctx, &chatMemberUpdated.From)
	default:
		log.Println("unknown chat member status:", chatMemberUpdated.NewChatMember.Status)
	}
}
