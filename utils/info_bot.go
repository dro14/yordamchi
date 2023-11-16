package utils

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func SendInfoMessage(text string) {
	config := tgbotapi.NewMessage(1331278972, "")
	config.ParseMode = tgbotapi.ModeMarkdownV2
	slices := Slice(text)
	for _, slice := range slices {
		config.Text = MarkdownV2(slice)
		_, err := bot.Request(config)
		if err != nil {
			log.Println("can't send info message:", err)
		}
	}
}

func SendLogFile(filepath string) {
	file := tgbotapi.FilePath(filepath)
	config := tgbotapi.NewDocument(1331278972, file)
	_, err := bot.Request(config)
	if err != nil {
		log.Println("can't send log file:", err)
	}
}
