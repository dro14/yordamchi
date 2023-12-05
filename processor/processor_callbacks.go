package processor

import (
	"context"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (p *Processor) newChatCallback(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	p.redis.DeleteContext(ctx)
	_, err := p.telegram.SendMessage(ctx, text.NewChat[lang(ctx)], 0, nil)
	if err != nil {
		log.Println("can't send new chat callback")
	}
	p.telegram.AnswerCallbackQuery(ctx, callbackQuery.ID, text.NewChatAnswer[lang(ctx)])
}

func (p *Processor) helpCallback(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	err := p.telegram.EditMessage(ctx, text.Help[lang(ctx)], callbackQuery.Message.MessageID, nil)
	if err != nil {
		log.Println("can't edit help callback")
	}
}

func (p *Processor) settingsCallback(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	switch callbackQuery.Data {
	case "settings", "settings1":
		err := p.telegram.EditMessage(ctx, text.Unlimited[lang(ctx)], callbackQuery.Message.MessageID, p.unlimitedButtons(ctx))
		if err != nil {
			log.Println("can't edit settings1 callback")
		}
	case "settings2":
		err := p.telegram.EditMessage(ctx, text.Premium[lang(ctx)], callbackQuery.Message.MessageID, p.premiumButtons(ctx))
		if err != nil {
			log.Println("can't edit settings2 callback")
		}
	}
}

func (p *Processor) languageCallback(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	ctx = context.WithValue(ctx, "language_code", callbackQuery.Data)
	p.redis.SetLang(ctx)
	p.redis.DeleteContext(ctx)
	p.telegram.SetCommands(ctx)
	p.postgres.SetLang(ctx, callbackQuery.From)
	p.start(ctx, callbackQuery.From)
	p.telegram.DeleteMessage(ctx, callbackQuery.Message.MessageID)
}

func (p *Processor) examplesCallback(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	err := p.telegram.EditMessage(ctx, text.Examples[lang(ctx)], callbackQuery.Message.MessageID, p.examplesButton(ctx))
	if err != nil {
		log.Println("can't edit examples callback")
	}
}

func (p *Processor) generateCallback(ctx context.Context, callbackQuery *tgbotapi.CallbackQuery) {
	messageID := callbackQuery.Message.MessageID
	ctx = context.WithValue(ctx, "style", callbackQuery.Data)

	err := p.telegram.EditMessage(ctx, text.Loading[lang(ctx)], messageID, nil)
	if err != nil {
		log.Println("can't send loading message")
		return
	}

	isTyping := p.telegram.SetTyping(ctx)
	defer isTyping.Store(false)

	prompt := p.redis.Generate(ctx)
	prompt = p.apis.Translate("auto", "en", prompt)
	path, caption := p.openai.ProcessGenerations(ctx, prompt)
	if path == "" {
		log.Println("can't process generations")
		err = p.telegram.EditMessage(ctx, caption, messageID, nil)
		if err != nil {
			log.Println("can't edit error message")
		}
		return
	}

	caption = p.apis.Translate("en", lang(ctx), caption)
	err = p.telegram.SendPhoto(ctx, path, "", caption, nil)
	if err != nil {
		log.Println("can't send photo")
		return
	}
	p.telegram.DeleteMessage(ctx, messageID)
	p.redis.DecrementImages(ctx)
}
