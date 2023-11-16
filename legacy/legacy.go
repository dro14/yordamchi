package legacy

import (
	"log"
	"os"

	"github.com/dro14/yordamchi/processor/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Legacy struct {
	bot *tgbotapi.BotAPI
}

func New() *Legacy {
	token, ok := os.LookupEnv("LEGACY_BOT_TOKEN")
	if !ok {
		log.Fatal("legacy bot token is not specified")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("can't initialize legacy bot:", err)
	}

	return &Legacy{
		bot: bot,
	}
}

func (l *Legacy) Redirect(message *tgbotapi.Message) {
	config := tgbotapi.NewMessage(message.From.ID, "")
	switch message.From.LanguageCode {
	case "uz", "":
		config.Text = text.LegacyMessage["uz"]
	case "ru":
		config.Text = text.LegacyMessage["ru"]
	default:
		config.Text = text.LegacyMessage["en"]
	}
	_, err := l.bot.Request(config)
	if err != nil {
		log.Println("can't send legacy message:", err)
	}
}
