package utils

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func SendInfoMessage(text, path string) {
	if path != "" {
		config := tgbotapi.NewPhoto(1331278972, tgbotapi.FilePath(path))
		config.Caption = text
		_, err := bot.Request(config)
		if err != nil {
			log.Println("can't send info message:", err)
		}
	} else if text != "" {
		config := tgbotapi.NewMessage(1331278972, "")
		config.ParseMode = tgbotapi.ModeMarkdownV2
		config.DisableWebPagePreview = true
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

func SendLogFiles() {
	sendLogFile("gin.log")
	sendLogFile("yordamchi.log")
}

func sendLogFile(filepath string) {
	file := tgbotapi.FilePath(filepath)
	config := tgbotapi.NewDocument(1331278972, file)
	_, err := bot.Request(config)
	if err != nil {
		log.Println("can't send log file:", err)
	}
}
