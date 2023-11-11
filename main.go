package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

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
		log.Fatalf("can't create telegram.log: %v", err)
	}

	time.Local, _ = time.LoadLocation("Asia/Tashkent")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(file)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		sig := <-sigChan
		log.Printf("Received %v, initiating shutdown...\n", sig)
		RecoverFromPanic()
		os.Exit(0)
	}()

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

func RecoverFromPanic() {
	if r := recover(); r != nil {
		file, _ := os.OpenFile("yordamchi.log", os.O_APPEND|os.O_WRONLY, 0644)
		_, _ = file.WriteString(r.(string) + "\n" + string(debug.Stack()))
		_ = file.Close()
	}
	info_bot.SendFile("yordamchi.log")
}

func stackTrace() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
