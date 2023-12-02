package processor

import (
	"context"
	"fmt"
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
	case "memory":
		p.memory(ctx)
	case "examples":
		p.examples(ctx)
	case "unlimited":
		p.unlimited(ctx)
	case "premium":
		p.premium(ctx)
	case "images":
		p.images(ctx)
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
	p.postgres.UserStarted(ctx, user)
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

func (p *Processor) memory(ctx context.Context) {

}

func (p *Processor) examples(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.Examples[lang(ctx)], 0, p.examplesButton(ctx))
	if err != nil {
		log.Println("can't send examples command")
	}
}

func (p *Processor) unlimited(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.Unlimited[lang(ctx)], 0, p.unlimitedButtons(ctx))
	if err != nil {
		log.Println("can't send unlimited command")
	}
}

func (p *Processor) premium(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.Premium[lang(ctx)], 0, p.premiumButtons(ctx))
	if err != nil {
		log.Println("can't send premium command")
	}
}

func (p *Processor) images(ctx context.Context) {
	userID := ctx.Value("user_id").(int64)
	config := tgbotapi.NewCopyMessage(userID, -1001924963699, 49)
	config.Caption = fmt.Sprintf(text.Image[lang(ctx)], p.redis.Images(ctx))
	config.ReplyMarkup = p.imageButtons(ctx)
	p.telegram.CopyMessage(ctx, &config)
}

func (p *Processor) generate(ctx context.Context, message *tgbotapi.Message) {
	if p.redis.Images(ctx) == 0 {
		p.images(ctx)
		return
	}
	generate, _ := strings.CutPrefix(message.Text, "/generate")
	p.redis.SetGenerate(ctx, strings.TrimSpace(generate))
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
