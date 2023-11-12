package telegram

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/dro14/yordamchi/client/telegram"
	"github.com/dro14/yordamchi/client/translator"
	"github.com/dro14/yordamchi/postgres"
	"github.com/dro14/yordamchi/processor/openai"
	"github.com/dro14/yordamchi/processor/telegram/button"
	"github.com/dro14/yordamchi/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func doCommand(ctx context.Context, message *tgbotapi.Message) bool {
	switch message.Command() {
	case "start":
		start(ctx, message)
	case "help":
		help(ctx)
	case "settings":
		settings(ctx)
	case "language":
		language(ctx)
	case "examples":
		examples(ctx)
	case "premium":
		premium(ctx)
	case "gpt4":
		gpt4(ctx)
	case "image":
		// TODO: add image command
	case "generate":
		generate(ctx, message)
	}
	return message.IsCommand()
}

func start(ctx context.Context, message *tgbotapi.Message) {
	_, err := telegram.Send(ctx, text.Start[lang(ctx)], 0, button.Start(lang(ctx)))
	if err != nil {
		log.Printf("can't send start command")
	}
	str, _ := strings.CutPrefix(message.Text, "/start ")
	joinedBy, _ := strconv.Atoi(str)
	postgres.JoinUser(ctx, message.From, int64(joinedBy))
}

func help(ctx context.Context) {
	_, err := telegram.Send(ctx, text.Help[lang(ctx)], 0, nil)
	if err != nil {
		log.Printf("can't send help command")
	}
}

func settings(ctx context.Context) {
	_, err := telegram.Send(ctx, msg(ctx), 0, button.Settings(ctx))
	if err != nil {
		log.Printf("can't send settings command")
	}
}

func language(ctx context.Context) {
	_, err := telegram.Send(ctx, text.Language, 0, button.Language())
	if err != nil {
		log.Printf("can't send language command")
	}
}

func examples(ctx context.Context) {
	_, err := telegram.Send(ctx, text.Examples[lang(ctx)], 0, button.Examples(lang(ctx)))
	if err != nil {
		log.Printf("can't send examples command")
	}
}

func premium(ctx context.Context) {
	_, err := telegram.Send(ctx, text.Premium[lang(ctx)], 0, button.Premium(ctx, lang(ctx)))
	if err != nil {
		log.Printf("can't send premium command")
	}
}

func gpt4(ctx context.Context) {
	_, err := telegram.Send(ctx, text.GPT4[lang(ctx)], 0, button.GPT4(ctx, lang(ctx)))
	if err != nil {
		log.Printf("can't send gpt4 command")
	}
}

func generate(ctx context.Context, message *tgbotapi.Message) {
	_, err := telegram.Send(ctx, text.Loading[lang(ctx)], message.MessageID, nil)
	if err != nil {
		log.Printf("can't send loading message")
		return
	}

	prompt := strings.ReplaceAll(message.Text, "/generate", "")
	prompt = strings.TrimSpace(prompt)
	prompt, _ = translator.Translate("auto", "en", prompt)

	photoURL, err := openai.ProcessGenerations(ctx, prompt)
	if err != nil {
		log.Printf("can't process generations: %s", err)
		return
	}

	telegram.SendPhoto(ctx, photoURL)
	telegram.Delete(ctx, message.MessageID)
}

func exhausted(ctx context.Context) {
	_, err := telegram.Send(ctx, text.Exhausted[lang(ctx)], 0, button.Premium(ctx, lang(ctx)))
	if err != nil {
		log.Printf("can't send exhausted message")
	}
}
