package processor

import (
	"context"
	"log"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) newChat(ctx context.Context) {
	p.redis.DeleteHistory(ctx)
	_, err := p.telegram.SendMessage(ctx, text.NewChat[lang(ctx)], 0, nil)
	if err != nil {
		log.Println("can't send new chat callback")
	}
}

func (p *Processor) examplesCallback(ctx context.Context, messageID int) {
	err := p.telegram.EditMessage(ctx, text.Examples[lang(ctx)], messageID, p.examplesButton(lang(ctx)))
	if err != nil {
		log.Println("can't edit examples callback")
	}
}

func (p *Processor) helpCallback(ctx context.Context, messageID int) {
	err := p.telegram.EditMessage(ctx, text.Help[lang(ctx)], messageID, nil)
	if err != nil {
		log.Println("can't edit help callback")
	}
}

func (p *Processor) model(ctx context.Context, messageID int, model string) {
	if model == p.redis.Model(ctx) {
		return
	} else if model == models.GPT4 {
		p.redis.GPT4(ctx)
	} else {
		p.redis.GPT3(ctx)
	}
	err := p.telegram.EditMessage(ctx, p.msg(ctx), messageID, p.settingsButtons(ctx))
	if err != nil {
		log.Println("can't edit model callback")
	}
}

func (p *Processor) languageCallback(ctx context.Context, message *tgbotapi.Message, lang string) {
	ctx = context.WithValue(ctx, "language_code", lang)
	p.redis.SetLang(ctx)
	p.redis.DeleteHistory(ctx)
	p.telegram.SetCommands(ctx)
	p.postgres.SetLang(ctx, message.From)
	p.start(ctx, message)
	p.telegram.DeleteMessage(ctx, message.MessageID)
}
