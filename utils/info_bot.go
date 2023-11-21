package utils

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func SendInfoMessage(text string, message *tgbotapi.Message) {
	if message != nil {
		if message.Photo != nil {
			photo := message.Photo[len(message.Photo)-1]
			file := tgbotapi.FileID(photo.FileID)
			config := tgbotapi.NewPhoto(1331278972, file)
			config.Caption = message.Caption
			_, err := bot.Request(config)
			if err != nil {
				log.Println("can't send info message:", err)
			}
		} else {
			text = message.Text
		}
	}
	if text != "" {
		config := tgbotapi.NewMessage(1331278972, "")
		config.ParseMode = tgbotapi.ModeMarkdownV2
		slices := Slice(text, 4096)
		for _, slice := range slices {
			config.Text = MarkdownV2(slice)
			_, err := bot.Request(config)
			if err != nil {
				log.Println("can't send info message:", err)
			}
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
