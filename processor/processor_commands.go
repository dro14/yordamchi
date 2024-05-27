package processor

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/dro14/yordamchi/storage/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) command(ctx context.Context, message *tgbotapi.Message) {
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
	case "system":
		p.system(ctx, message)
	case "files":
		p.files(ctx, message)
	}
}

func (p *Processor) start(ctx context.Context, user *tgbotapi.User) {
	p.postgres.UserStarted(ctx, user)
	name := p.postgres.User(ctx, user)
	if name == "" {
		name = text.DearUser[lang(ctx)]
	}
	Text := fmt.Sprintf(text.Start[lang(ctx)], name)
	_, err := p.telegram.SendMessage(ctx, Text, 0, p.startButton(ctx))
	if err != nil {
		log.Println("can't send start command")
	}
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
	switch userStatus(ctx) {
	case redis.StatusFree, redis.StatusExhausted, redis.StatusUnknown:
		p.paidFeature(ctx)
		return
	default:
	}

	var Text string
	system, found := strings.CutPrefix(p.redis.System(ctx), "USER: ")
	if found {
		Text = fmt.Sprintf(text.MemorySystem[lang(ctx)], system)
	} else {
		Text = text.MemoryEmpty[lang(ctx)]
	}
	Text = fmt.Sprintf(text.Memory[lang(ctx)], Text, p.service.Memory(ctx))

	_, err := p.telegram.SendMessage(ctx, Text, 0, p.newChatButton(ctx))
	if err != nil {
		log.Println("can't send memory command")
	}
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
	caption := fmt.Sprintf(text.Image[lang(ctx)], p.redis.Images(ctx))
	err := p.telegram.SendPhoto(ctx, "images.jpeg", caption, p.imageButtons(ctx))
	if err != nil {
		log.Println("can't send images command")
	}
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

func (p *Processor) system(ctx context.Context, message *tgbotapi.Message) {
	switch userStatus(ctx) {
	case redis.StatusFree, redis.StatusExhausted, redis.StatusUnknown:
		p.paidFeature(ctx)
		return
	default:
	}

	system, _ := strings.CutPrefix(message.Text, "/system")
	if model(ctx) == models.GPT3 && lang(ctx) == "uz" {
		system = p.apis.Translate("auto", "en", system)
	}
	system = strings.TrimSpace(system)
	p.redis.SetSystem(ctx, system)
	Text := fmt.Sprintf(text.System[lang(ctx)], system)

	_, err := p.telegram.SendMessage(ctx, Text, 0, nil)
	if err != nil {
		log.Println("can't send system command")
	}
}

func (p *Processor) files(ctx context.Context, message *tgbotapi.Message) {
	if message.From.ID == 1331278972 {
		_, err := p.telegram.SendMessage(ctx, p.service.Files(ctx), 0, nil)
		if err != nil {
			log.Println("can't send files command")
		}
	}
}
