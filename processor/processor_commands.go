package processor

import (
	"context"
	"log"
	"strings"

	"github.com/dro14/yordamchi/processor/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) doCommand(ctx context.Context, message *tgbotapi.Message) bool {
	switch message.Command() {
	case "start":
		p.start(ctx, message.From)
	case "help":
		p.help(ctx)
	case "settings":
		p.settings(ctx)
	case "language":
		p.language(ctx)
	case "examples":
		p.examples(ctx)
	case "premium":
		p.premium(ctx)
	case "image":
		// TODO: add image command
	case "generate":
		p.generate(ctx, message)
	case "logs":
		p.logs(ctx, message)
	}
	return message.IsCommand()
}

func (p *Processor) start(ctx context.Context, user *tgbotapi.User) {
	_, err := p.telegram.SendMessage(ctx, text.Start[lang(ctx)], 0, p.startButton(ctx))
	if err != nil {
		log.Println("can't send start command")
	}
	p.postgres.JoinUser(ctx, user)
}

func (p *Processor) help(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.Help[lang(ctx)], 0, nil)
	if err != nil {
		log.Println("can't send help command")
	}
}

func (p *Processor) settings(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, p.msg(ctx), 0, p.settingsButton(ctx))
	if err != nil {
		log.Println("can't send settings command")
	}
}

func (p *Processor) language(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.Language, 0, p.languageButtons())
	if err != nil {
		log.Println("can't send language command")
	}
}

func (p *Processor) examples(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.Examples[lang(ctx)], 0, p.examplesButton(ctx))
	if err != nil {
		log.Println("can't send examples command")
	}
}

func (p *Processor) premium(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.Premium[lang(ctx)], 0, p.premiumButtons(ctx))
	if err != nil {
		log.Println("can't send premium command")
	}
}

func (p *Processor) generate(ctx context.Context, message *tgbotapi.Message) {
	messageID, err := p.telegram.SendMessage(ctx, text.Loading[lang(ctx)], message.MessageID, nil)
	if err != nil {
		log.Println("can't send loading message")
		return
	}

	prompt := strings.ReplaceAll(message.Text, "/generate", "")
	prompt = strings.TrimSpace(prompt)
	prompt = p.apis.Translate("auto", "en", prompt)

	photoURL, err := p.openai.ProcessGenerations(ctx, prompt)
	if err != nil {
		log.Println("can't process generations:", err)
		return
	}

	p.telegram.SendPhoto(ctx, photoURL)
	p.telegram.DeleteMessage(ctx, messageID)
}

func (p *Processor) logs(ctx context.Context, message *tgbotapi.Message) {
	if message.From.ID == 1331278972 {
		p.telegram.SendFile(ctx, "gin.log")
		p.telegram.SendFile(ctx, "yordamchi.log")
	}
}
