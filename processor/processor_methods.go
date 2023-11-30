package processor

import (
	"context"
	"errors"
	"log"

	"github.com/dro14/yordamchi/clients/service"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) exhausted(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.Exhausted[lang(ctx)], 0, p.settingsButton(ctx))
	if err != nil {
		log.Println("can't send exhausted message")
	}
}

func (p *Processor) paidFeature(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.PaidFeature[lang(ctx)], 0, p.settingsButton(ctx))
	if err != nil {
		log.Println("can't send premium feature message")
	}
}

func (p *Processor) processFile(ctx context.Context, message *tgbotapi.Message) {
	messageID, err := p.telegram.SendMessage(ctx, text.Loading[lang(ctx)], message.MessageID, nil)
	if err != nil {
		log.Println("can't send loading message")
		return
	}

	err = p.service.Load(ctx, message.Document)
	if errors.Is(err, service.ErrUnsupportedFormat) {
		err = p.telegram.EditMessage(ctx, text.UnsupportedFormat[lang(ctx)], messageID, nil)
		if err != nil {
			log.Println("can't edit message")
		}
	} else if err != nil {
		log.Println("can't load file:", err)
		err = p.telegram.EditMessage(ctx, text.RequestFailed[lang(ctx)], messageID, nil)
		if err != nil {
			log.Println("can't edit message")
		}
	} else {
		err = p.telegram.EditMessage(ctx, text.Loaded[lang(ctx)], messageID, nil)
		if err != nil {
			log.Println("can't edit message")
		}
	}
}
