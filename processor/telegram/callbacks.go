package telegram

import (
	"context"
	"log"

	"github.com/dro14/yordamchi/client/telegram"
	"github.com/dro14/yordamchi/processor/telegram/button"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
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

func premiumCallback(ctx context.Context, messageID int) {

	err := telegram.EditMessage(ctx, text.Premium[lang(ctx)], messageID, button.Premium(ctx, lang(ctx)))
	if err != nil {
		log.Printf("can't edit premium callback")
	}
}

func modelCallback(ctx context.Context, messageID int, model string) {

	var err error
	if model == redis.Model(ctx) {
		return
	} else if model == "gpt-4" {
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

func translatorEnabled(ctx context.Context, messageID int) {

	redis.SetLang(ctx, "uz")
	_, err := telegram.SendMessage(ctx, text.TranslatorEnabled[lang(ctx)], 0, nil)
	if err != nil {
		log.Printf("can't send translate callback")
	}

	err = telegram.EditMessage(ctx, text.Translate[lang(ctx)], messageID, nil)
	if err != nil {
		log.Printf("can't edit translate callback")
	}
}

func translatorDisabled(ctx context.Context, messageID int) {

	redis.SetLang(ctx, "-")
	_, err := telegram.SendMessage(ctx, text.TranslatorDisabled[lang(ctx)], 0, nil)
	if err != nil {
		log.Printf("can't send translate callback")
	}

	err = telegram.EditMessage(ctx, text.Translate[lang(ctx)], messageID, nil)
	if err != nil {
		log.Printf("can't edit translate callback")
	}
}
