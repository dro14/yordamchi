package processor

import (
	"context"
	"log"

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

}
