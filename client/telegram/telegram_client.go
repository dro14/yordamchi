package telegram

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func Init() {
	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatalf("bot token is not specified")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("can't initialize bot: %v", err)
	}
}
