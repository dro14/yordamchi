package telegram

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/dro14/yordamchi/client/bobdev"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/processor/telegram/button"
	"github.com/dro14/yordamchi/processor/telegram/text"
	"github.com/gotd/td/tg"
)

func (p *Processor) newChatCallback(ctx context.Context) {

	p.Cache.DeleteContext(ctx)

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

	err := p.Client.EditMessage(ctx, text.Premium[lang(ctx)], messageID, button.Premium(lang(ctx)))
	if err != nil {
		log.Printf("can't edit premium callback")
	}
}

func (p *Processor) confirmCallback(ctx context.Context, data string, user *tg.User) {

	before, after, found := strings.Cut(data, ":")
	if !found {
		log.Printf("can't parse callback data: %q", data)
		return
	}
	price, _ := strconv.Atoi(before)
	requests, _ := strconv.Atoi(after)

	request := &types.Request{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		LanguageCode: lang(ctx),
		Price:        price,
		Requests:     requests,
	}

	err := p.Client.EditMessage(ctx, text.Confirm[lang(ctx)], ctx.Value("message_id").(int), button.URLButton("Payme", bobdev.Payme(request)))
	if err != nil {
		log.Printf("can't edit confirm callback")
	}
}
