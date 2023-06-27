package telegram

import (
	"context"
	"github.com/dro14/yordamchi/client/openai"
	"github.com/dro14/yordamchi/client/telegram"
	"github.com/dro14/yordamchi/payme"
	"github.com/dro14/yordamchi/postgres"
	"github.com/dro14/yordamchi/processor/telegram/legacy_bot"
	"log"
	"net/http"
	"time"

	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gotd/td/tg"
)

func Init() {

	time.Local, _ = time.LoadLocation("Asia/Tashkent")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	telegram.Init()
	openai.Init()
	postgres.Init()
	redis.Init()
	legacy_bot.Init()
	payme.Run()
}

func ProcessUpdate(c *gin.Context) {

	update := &tgbotapi.Update{}
	if err := c.ShouldBindJSON(update); err != nil {
		log.Printf("can't bind json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch {
	case update.Message != nil:
		ProcessMessage()
	case update.CallbackQuery != nil:
		ProcessCallbackQuery()
	case update.MyChatMember != nil:

	}

}

func ProcessMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) error {

	ctx, message, user := messageUpdate(ctx, entities, update)
	if message == nil && user == nil {
		return nil
	}
	defer blockedUsers.Delete(user.ID)

	done := doCommand(ctx, message, user)
	if done {
		return nil
	}

	switch redis.UserStatus(ctx) {
	case types.UnknownStatus:
		log.Printf("unknown user status: %d", user.ID)
	case types.BlockedStatus:
		blocked(ctx)
	case types.GPT4Status:
		if redis.GPT4Tokens(ctx) > 0 {
			Stream(ctx, message, user, "gpt-4")
		} else {
			gpt4(ctx)
		}
	case types.PremiumStatus:
		Stream(ctx, message, user, "true")
	case types.FreeStatus:
		Stream(ctx, message, user, "false")
	case types.ExhaustedStatus:
		exhausted(ctx)
	}

	return nil
}

func ProcessCallbackQuery(ctx context.Context, entities tg.Entities, update *tg.UpdateBotCallbackQuery) error {

	ctx, data := callbackUpdate(ctx, entities, update)

	switch data {
	case "new_chat":
		p.newChatCallback(ctx)
	case "examples":
		p.examplesCallback(ctx, update.MsgID)
	case "help":
		p.helpCallback(ctx, update.MsgID)
	case "premium":
		p.premiumCallback(ctx, update.MsgID)
	case "gpt-3.5-turbo", "gpt-4":
		p.modelCallback(ctx, update.MsgID, data)
	default:
		log.Printf("unknown callback data: %v", data)
	}

	return nil
}

func ProcessBotStopped(ctx context.Context, entities tg.Entities, update *tg.UpdateBotStopped) error {

	ctx, user := botStoppedUpdate(ctx, entities, update)

	if update.Stopped {
		p.deactivated(ctx, user)
	} else {
		p.rejoined(ctx, user)
	}

	return nil
}
