package processor

import (
	"context"
	"log"

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
