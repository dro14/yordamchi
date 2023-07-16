package info_bot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

var bot *tgbotapi.BotAPI

func Init() {

	token, ok := os.LookupEnv("INFO_BOT_TOKEN")
	if !ok {
		log.Fatalf("info bot token is not specified")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("can't initialize info bot: %v", err)
	}
}

func Send(text string) {

	config := tgbotapi.NewMessage(1331278972, text)
	_, err := bot.Send(config)
	if err != nil {
		log.Printf("can't send info message: %v", err)
	}
}
