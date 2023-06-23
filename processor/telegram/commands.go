package telegram

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/postgres"
	"github.com/dro14/yordamchi/processor/telegram/button"
	"github.com/dro14/yordamchi/text"
	"github.com/gotd/td/tg"
)

func (p *Processor) doCommand(ctx context.Context, message *tg.Message, user *tg.User) bool {

	switch command(message) {
	case "start":
		p.start(ctx, message, user)
	case "help":
		p.help(ctx)
	case "settings":
		p.settings(ctx)
	case "examples":
		p.examples(ctx)
	case "premium":
		p.premium(ctx)
	case "donate":
		p.donate(ctx)
	default:
		return false
	}

	return true
}

func (p *Processor) start(ctx context.Context, message *tg.Message, user *tg.User) {

	_, err := p.Client.SendMessage(ctx, text.Start[lang(ctx)], 0, button.Start(lang(ctx)))
	if err != nil {
		log.Printf("can't send start command")
	}

	str, _ := strings.CutPrefix(message.Message, "/start ")
	joinedBy, _ := strconv.Atoi(str)
	postgres.JoinUser(ctx, user, int64(joinedBy))
}

func (p *Processor) help(ctx context.Context) {

	_, err := p.Client.SendMessage(ctx, text.Help[lang(ctx)], 0, nil)
	if err != nil {
		log.Printf("can't send help command")
	}
}

func (p *Processor) settings(ctx context.Context) {

	message := settingsMessage(ctx, lang(ctx))
	if len(message) == 0 {
		log.Printf("can't get user settings")
		return
	}

	_, err := p.Client.SendMessage(ctx, message, 0, button.Settings(lang(ctx)))
	if err != nil {
		log.Printf("can't send settings command")
	}
}

func (p *Processor) examples(ctx context.Context) {

	_, err := p.Client.SendMessage(ctx, text.Examples[lang(ctx)], 0, button.Examples(lang(ctx)))
	if err != nil {
		log.Printf("can't send examples command")
	}
}

func (p *Processor) premium(ctx context.Context) {

	_, err := p.Client.SendMessage(ctx, text.Premium[lang(ctx)], 0, button.Premium(lang(ctx)))
	if err != nil {
		log.Printf("can't send premium command")
	}
}

func (p *Processor) donate(ctx context.Context) {

	_, err := p.Client.SendMessage(ctx, text.Donate[lang(ctx)], 0, button.URLButton("Payme", constants.DonationURL))
	if err != nil {
		log.Printf("can't send donate command")
	}
}
