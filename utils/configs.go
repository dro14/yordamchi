package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SetConfigs() {
	time.Local, _ = time.LoadLocation("Asia/Tashkent")

	file, err := os.Create("yordamchi.log")
	if err != nil {
		log.Fatal("can't create yordamchi.log:", err)
	}
	log.SetOutput(file)
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

	go keepServiceAlive()
}

func keepServiceAlive() {
	for {
		resp, err := http.Get("https://yordamchi-service.victoriousriver-fffd2d70.westeurope.azurecontainerapps.io")
		if err != nil {
			log.Println("can't ping service:", err)
			continue
		}

		response := make(map[string]string)
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			log.Println("can't decode response:", err)
			continue
		}

		if response["message"] != "Hello, Yordamchi!" {
			log.Println("Yordamchi service is not alive!")
		}
		time.Sleep(1 * time.Minute)
	}
}
