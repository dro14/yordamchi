package utils

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func SendInfoMessage(fromChatID int64, messageID int) {
	if fromChatID == 1792604195 || fromChatID == -1001924963699 {
		config := tgbotapi.NewCopyMessage(1331278972, fromChatID, messageID)
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
