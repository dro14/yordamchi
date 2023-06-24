package telegram

import (
	"context"
	"github.com/dro14/yordamchi/processor/telegram/button"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
	"log"
)

func (p *Processor) newChatCallback(ctx context.Context) {

	redis.DeleteContext(ctx)
	_, err := p.Client.SendMessage(ctx, text.NewChat[lang(ctx)], 0, nil)
	if err != nil {
		log.Printf("can't send new chat callback")
	}
}

func (p *Processor) examplesCallback(ctx context.Context, messageID int) {

	err := p.Client.EditMessage(ctx, text.Examples[lang(ctx)], messageID, button.Examples(lang(ctx)))
	if err != nil {
		log.Printf("can't edit examples callback")
	}
}

func (p *Processor) helpCallback(ctx context.Context, messageID int) {

	err := p.Client.EditMessage(ctx, text.Help[lang(ctx)], messageID, nil)
	if err != nil {
		log.Printf("can't edit help callback")
	}
}

func (p *Processor) premiumCallback(ctx context.Context, messageID int) {

	err := p.Client.EditMessage(ctx, text.Premium[lang(ctx)], messageID, button.Premium(ctx, lang(ctx)))
	if err != nil {
		log.Printf("can't edit premium callback")
	}
}
