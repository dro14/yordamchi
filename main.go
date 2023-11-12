package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dro14/yordamchi/lib/recover"
	"github.com/dro14/yordamchi/payme"
	"github.com/dro14/yordamchi/processor/telegram"
	"github.com/dro14/yordamchi/processor/telegram/info_bot"
	"github.com/dro14/yordamchi/processor/telegram/legacy_bot"
	"github.com/dro14/yordamchi/recovery"
	"github.com/gin-gonic/gin"
)

func main() {
	file, err := os.Create("yordamchi.log")
	if err != nil {
		log.Fatalf("can't create yordamchi.log: %s", err)
	}

	time.Local, _ = time.LoadLocation("Asia/Tashkent")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(file)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer recover.Recover(sigChan)

	file, err = os.Create("gin.log")
	if err != nil {
		log.Fatalf("can't create gin.log: %v", err)
	}
	gin.DefaultWriter = file
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.POST("/main", telegram.ProcessUpdate)
	r.POST("/legacy", legacy_bot.Reply)
	r.POST("/payme", payme.Handler)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	telegram.Init()
	recovery.Init()
	info_bot.SendMessage("@yordamchi_ai_bot restarted")

	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("can't run server: %v", err)
	}
}
