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
		p.image(ctx)
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

func (p *Processor) image(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.Image[lang(ctx)], 0, p.imageButtons(ctx))
	if err != nil {
		log.Println("can't send image command")
	}
}

func (p *Processor) generate(ctx context.Context, message *tgbotapi.Message) {
	prompt := strings.ReplaceAll(message.Text, "/generate", "")
	prompt = strings.TrimSpace(prompt)
	p.redis.StorePrompt(ctx, prompt)
	_, err := p.telegram.SendMessage(ctx, text.Generate[lang(ctx)], message.MessageID, p.generateButtons(ctx))
	if err != nil {
		log.Println("can't send generate command")
	}
}

func (p *Processor) logs(ctx context.Context, message *tgbotapi.Message) {
	if message.From.ID == 1331278972 {
		p.telegram.SendFile(ctx, "gin.log")
		p.telegram.SendFile(ctx, "yordamchi.log")
	}
}
