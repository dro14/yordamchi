package legacy

import (
	"github.com/dro14/yordamchi/utils"
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

func (l *Legacy) Redirect(update *tgbotapi.Update) {
	if update.Message == nil {
		log.Printf("unknown update type:\n%+v", update)
		return
	}

	user := update.Message.From
	config := tgbotapi.NewMessage(user.ID, "")
	switch user.LanguageCode {
	case "uz", "":
		config.Text = text.LegacyMessage["uz"]
	case "ru":
		config.Text = text.LegacyMessage["ru"]
	default:
		config.Text = text.LegacyMessage["en"]
	}
	config.Text = utils.MarkdownV2(config.Text)
	config.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := l.bot.Request(config)
	if err != nil {
		log.Println("can't send legacy message:", err)
	}
}
