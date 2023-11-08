package telegram

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/dro14/yordamchi/client/telegram"
	"github.com/dro14/yordamchi/postgres"
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
		// TODO: add generate command
	}
	return message.IsCommand()
}

func start(ctx context.Context, message *tgbotapi.Message) {
	_, err := telegram.SendMessage(ctx, text.Start[lang(ctx)], 0, button.Start(lang(ctx)))
	if err != nil {
		log.Printf("can't send start command")
	}
	str, _ := strings.CutPrefix(message.Text, "/start ")
	joinedBy, _ := strconv.Atoi(str)
	postgres.JoinUser(ctx, message.From, int64(joinedBy))
}

func help(ctx context.Context) {
	_, err := telegram.SendMessage(ctx, text.Help[lang(ctx)], 0, nil)
	if err != nil {
		log.Printf("can't send help command")
	}
}

func settings(ctx context.Context) {
	_, err := telegram.SendMessage(ctx, msg(ctx, lang(ctx)), 0, button.Settings(ctx))
	if err != nil {
		log.Printf("can't send settings command")
	}
}

func language(ctx context.Context) {
	_, err := telegram.SendMessage(ctx, text.Language, 0, button.Language())
	if err != nil {
		log.Printf("can't send language command")
	}
}

func examples(ctx context.Context) {
	_, err := telegram.SendMessage(ctx, text.Examples[lang(ctx)], 0, button.Examples(lang(ctx)))
	if err != nil {
		log.Printf("can't send examples command")
	}
}

func premium(ctx context.Context) {
	_, err := telegram.SendMessage(ctx, text.Premium[lang(ctx)], 0, button.Premium(ctx, lang(ctx)))
	if err != nil {
		log.Printf("can't send premium command")
	}
}

func gpt4(ctx context.Context) {
	_, err := telegram.SendMessage(ctx, text.GPT4[lang(ctx)], 0, button.GPT4(ctx, lang(ctx)))
	if err != nil {
		log.Printf("can't send gpt4 command")
	}
}

func exhausted(ctx context.Context) {
	_, err := telegram.SendMessage(ctx, text.Exhausted[lang(ctx)], 0, button.Premium(ctx, lang(ctx)))
	if err != nil {
		log.Printf("can't send exhausted message")
	}
}
