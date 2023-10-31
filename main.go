package main

import (
	"log"
	"os"
	"time"

	"github.com/dro14/yordamchi/payme"
	"github.com/dro14/yordamchi/processor/telegram"
	"github.com/dro14/yordamchi/processor/telegram/info_bot"
	"github.com/dro14/yordamchi/processor/telegram/legacy_bot"
	"github.com/dro14/yordamchi/recovery"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {

	time.Local, _ = time.LoadLocation("Asia/Tashkent")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	gin.SetMode(gin.ReleaseMode)

	file, err := os.Create("gin.log")
	if err != nil {
		log.Fatalf("can't create gin.log: %v", err)
	}
	gin.DefaultWriter = file

	r := gin.Default()
	r.POST("/main", telegram.ProcessUpdate)
	r.POST("/legacy", legacy_bot.Reply)
	r.POST("/payme", payme.Handler)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "80"
	}

	telegram.Init()
	recovery.Init()
	info_bot.Send("@yordamchi_ai_bot restarted")

	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("can't run server: %v", err)
	}
}
