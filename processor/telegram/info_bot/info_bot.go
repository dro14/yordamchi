package info_bot

import (
	"log"
	"os"

	"github.com/dro14/yordamchi/lib/functions"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func SendMessage(text string) {
	config := tgbotapi.NewMessage(1331278972, "")
	config.ParseMode = tgbotapi.ModeMarkdownV2
	slices := functions.Slice(text)
	for _, slice := range slices {
		config.Text = functions.MarkdownV2(slice)
		_, err := bot.Send(config)
		if err != nil {
			log.Printf("can't send info message: %s", err)
		}
	}
}

func SendFile(path string) {
	file := tgbotapi.FilePath(path)
	config := tgbotapi.NewDocument(1331278972, file)
	_, err := bot.Request(config)
	if err != nil {
		log.Printf("can't send log file: %s", err)
	}
}
