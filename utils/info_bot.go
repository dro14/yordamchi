package utils

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func SendInfoMessage(text string) {
	config := tgbotapi.NewMessage(AdminUserID, text)
	_, err := bot.Request(config)
	if err != nil {
		log.Println("can't send info message:", err)
	}
}

func SendLogFiles() {
	sendLogFile("gin.log")
	sendLogFile("yordamchi.log")
}

func sendLogFile(filepath string) {
	file := tgbotapi.FilePath(filepath)
	config := tgbotapi.NewDocument(AdminUserID, file)
	_, err := bot.Request(config)
	if err != nil {
		log.Println("can't send log file:", err)
	}
}
