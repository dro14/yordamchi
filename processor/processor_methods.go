package processor

import (
	"context"
	"log"

	"github.com/dro14/yordamchi/processor/text"
)

func (p *Processor) exhausted(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.Exhausted[lang(ctx)], 0, p.settingsButton(ctx))
	if err != nil {
		log.Println("can't send exhausted message")
	}
}

func (p *Processor) premiumFeature(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.PremiumFeature[lang(ctx)], 0, p.settingsButton(ctx))
	if err != nil {
		log.Println("can't send premium feature message")
	}
}
