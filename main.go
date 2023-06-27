package main

import (
	"github.com/dro14/yordamchi/processor/telegram/legacy_bot"
	"log"
	"os"

	"github.com/dro14/yordamchi/payme"
	"github.com/dro14/yordamchi/processor/telegram"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {

	telegram.Init()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	file, err := os.Create("gin.log")
	if err != nil {
		log.Fatalf("can't create gin.log: %v", err)
	}

	r := gin.Default()
	gin.DefaultWriter = file
	gin.SetMode(gin.ReleaseMode)

	r.POST("/main", telegram.ProcessUpdate)
	r.POST("/legacy", legacy_bot.Reply)
	r.POST("/payme", payme.Handler)

	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("can't run server: %v", err)
	}
}
