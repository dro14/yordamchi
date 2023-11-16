package telegram

import (
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	token string
	bot   *tgbotapi.BotAPI
}

func New() *Telegram {
	token, ok := os.LookupEnv("MAIN_BOT_TOKEN")
	if !ok {
		log.Fatal("main bot token is not specified")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("can't initialize bot:", err)
	}

	return &Telegram{
		token: token,
		bot:   bot,
	}
}
