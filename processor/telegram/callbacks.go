package telegram

import (
	"context"
	"log"

	"github.com/dro14/yordamchi/client/telegram"
	"github.com/dro14/yordamchi/lib/models"
	"github.com/dro14/yordamchi/processor/telegram/button"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func newChatCallback(ctx context.Context) {
	redis.DeleteContext(ctx)
	_, err := telegram.SendMessage(ctx, text.NewChat[lang(ctx)], 0, nil)
	if err != nil {
		log.Printf("can't send new chat callback")
	}
}

func examplesCallback(ctx context.Context, messageID int) {
	err := telegram.EditMessage(ctx, text.Examples[lang(ctx)], messageID, button.Examples(lang(ctx)))
	if err != nil {
		log.Printf("can't edit examples callback")
	}
}

func helpCallback(ctx context.Context, messageID int) {
	err := telegram.EditMessage(ctx, text.Help[lang(ctx)], messageID, nil)
	if err != nil {
		log.Printf("can't edit help callback")
	}
}

func modelCallback(ctx context.Context, messageID int, model string) {
	var err error
	if model == redis.Model(ctx) {
		return
	} else if model == models.GPT4 {
		err = redis.GPT4(ctx)
	} else {
		err = redis.GPT3(ctx)
	}
	if err != nil {
		log.Printf("can't set model callback: %v", err)
		return
	}

	err = telegram.EditMessage(ctx, msg(ctx, lang(ctx)), messageID, button.Settings(ctx))
	if err != nil {
		log.Printf("can't edit model callback")
	}
}

func languageChosenCallback(ctx context.Context, message *tgbotapi.Message, lang string) {
	ctx = context.WithValue(ctx, "language_code", lang)
	redis.SetLang(ctx)
	redis.DeleteContext(ctx)
	telegram.SetCommands(ctx)
	start(ctx, message)

	err := telegram.Delete(ctx, message.MessageID)
	if err != nil {
		log.Printf("can't delete language command message")
	}
}
