package utils

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SetConfigs() {
	file, err := os.Create("yordamchi.log")
	if err != nil {
		log.Fatal("can't create yordamchi.log:", err)
	}
	// log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	file, err = os.Create("gin.log")
	if err != nil {
		log.Fatal("can't create gin.log:", err)
	}
	gin.DefaultWriter = file
	gin.SetMode(gin.ReleaseMode)

	token, ok := os.LookupEnv("INFO_BOT_TOKEN")
	if !ok {
		log.Fatal("info bot token is not specified")
	}

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("can't initialize info bot:", err)
	}
}
