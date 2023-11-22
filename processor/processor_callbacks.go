package processor

import (
	"context"
	"github.com/dro14/yordamchi/utils"
	"log"
	"strings"

	"github.com/dro14/yordamchi/processor/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) newChatCallback(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	p.redis.DeleteHistory(ctx)
	_, err := p.telegram.SendMessage(ctx, text.NewChat[lang(ctx)], 0, nil)
	if err != nil {
		log.Println("can't send new chat callback")
	}
	p.telegram.AnswerCallbackQuery(ctx, callbackQuery.ID, text.NewChatAnswer[lang(ctx)])
}

func (p *Processor) helpCallback(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	err := p.telegram.EditMessage(ctx, text.Help[lang(ctx)], callbackQuery.Message.MessageID, nil)
	if err != nil {
		log.Println("can't edit help callback")
	}
}

func (p *Processor) settingsCallback(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	err := p.telegram.EditMessage(ctx, text.Premium[lang(ctx)], callbackQuery.Message.MessageID, p.premiumButtons(ctx))
	if err != nil {
		log.Println("can't edit settings callback")
	}
}

func (p *Processor) languageCallback(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	ctx = context.WithValue(ctx, "language_code", callbackQuery.Data)
	p.redis.SetLang(ctx)
	p.redis.DeleteHistory(ctx)
	p.telegram.SetCommands(ctx)
	p.postgres.SetLang(ctx, callbackQuery.From)
	p.start(ctx, callbackQuery.From)
	p.telegram.DeleteMessage(ctx, callbackQuery.Message.MessageID)
}

func (p *Processor) examplesCallback(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	err := p.telegram.EditMessage(ctx, text.Examples[lang(ctx)], callbackQuery.Message.MessageID, p.examplesButton(ctx))
	if err != nil {
		log.Println("can't edit examples callback")
	}
}

func (p *Processor) generateCallback(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	style, prompt, found := strings.Cut(callbackQuery.Data, utils.Delimiter)
	if !found {
		log.Println("unknown callback data:", callbackQuery.Data)
		return
	}

	messageID := callbackQuery.Message.MessageID
	ctx = context.WithValue(ctx, "style", style)

	err := p.telegram.EditMessage(ctx, text.Loading[lang(ctx)], messageID, nil)
	if err != nil {
		log.Println("can't send loading message")
		return
	}

	photoURL, revisedPrompt := p.openai.ProcessGenerations(ctx, prompt)
	if photoURL == "" {
		log.Println("can't process generations:", err)
		err = p.telegram.EditMessage(ctx, revisedPrompt, messageID, nil)
		if err != nil {
			log.Println("can't edit error message")
		}
		return
	}

	revisedPrompt = p.apis.Translate("en", lang(ctx), revisedPrompt)
	p.telegram.SendPhoto(ctx, photoURL, revisedPrompt)
	p.telegram.DeleteMessage(ctx, messageID)
}
