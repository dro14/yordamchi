package telegram

import (
	"context"
	"log"

	"github.com/dro14/yordamchi/client/ocr"
	"github.com/dro14/yordamchi/client/openai"
	"github.com/dro14/yordamchi/client/telegram"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/payme"
	"github.com/dro14/yordamchi/postgres"
	"github.com/dro14/yordamchi/processor/telegram/info_bot"
	"github.com/dro14/yordamchi/processor/telegram/legacy_bot"
	"github.com/dro14/yordamchi/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init() {
	redis.Init()
	telegram.Init()
	openai.Init()
	postgres.Init()
	info_bot.Init()
	legacy_bot.Init()
	payme.Init()
	ocr.Init()
}

func ProcessUpdate(c *gin.Context) {

	update := &tgbotapi.Update{}
	if err := c.ShouldBindJSON(update); err != nil {
		log.Printf("can't bind json: %v", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	switch {
	case update.Message != nil:
		go ProcessMessage(context.Background(), update.Message)
	case update.CallbackQuery != nil:
		go ProcessCallbackQuery(context.Background(), update.CallbackQuery)
	case update.MyChatMember != nil:
		go ProcessMyChatMember(context.Background(), update.MyChatMember)
	default:
		log.Printf("unknown update type:\n%+v", update)
	}

	c.JSON(200, gin.H{"ok": true})
}

func ProcessMessage(ctx context.Context, message *tgbotapi.Message) {

	ctx, allow, err := messageUpdate(ctx, message)
	if !allow {
		return
	}
	defer blockedUsers.Delete(message.From.ID)

	done := doCommand(ctx, message)
	if err != nil {
		language(ctx)
		return
	} else if done {
		return
	}

	switch redis.UserStatus(ctx) {
	case types.GPT4Status:
		if redis.GPT4Tokens(ctx) > 0 {
			Stream(ctx, message, "gpt-4")
		} else {
			gpt4(ctx)
		}
	case types.PremiumStatus:
		Stream(ctx, message, "true")
	case types.FreeStatus:
		Stream(ctx, message, "false")
	case types.ExhaustedStatus:
		exhausted(ctx)
	default:
		log.Printf("unknown user status: %d", message.From.ID)
	}
}

func ProcessCallbackQuery(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {

	ctx = context.WithValue(ctx, "date", callbackQuery.Message.Date)
	ctx = context.WithValue(ctx, "user_id", callbackQuery.From.ID)
	ctx, _ = redis.Lang(ctx, callbackQuery.From.LanguageCode)

	switch callbackQuery.Data {
	case "new_chat":
		newChatCallback(ctx)
	case "examples":
		examplesCallback(ctx, callbackQuery.Message.MessageID)
	case "help":
		helpCallback(ctx, callbackQuery.Message.MessageID)
	case "gpt-3.5-turbo", "gpt-4":
		modelCallback(ctx, callbackQuery.Message.MessageID, callbackQuery.Data)
	case "uz", "ru", "en":
		languageChosen(ctx, callbackQuery.Message.MessageID, callbackQuery.Data)
	default:
		log.Printf("unknown callback data: %v", callbackQuery.Data)
	}
}

func ProcessMyChatMember(ctx context.Context, chatMemberUpdated *tgbotapi.ChatMemberUpdated) {

	ctx = context.WithValue(ctx, "date", chatMemberUpdated.Date)
	ctx = context.WithValue(ctx, "user_id", chatMemberUpdated.From.ID)
	ctx, _ = redis.Lang(ctx, chatMemberUpdated.From.LanguageCode)

	switch chatMemberUpdated.NewChatMember.Status {
	case "kicked":
		postgres.DeactivateUser(ctx, &chatMemberUpdated.From)
	case "member":
		postgres.RejoinUser(ctx, &chatMemberUpdated.From)
	default:
		log.Printf("unknown chat member status: %v", chatMemberUpdated.NewChatMember.Status)
	}
}
