package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dro14/yordamchi/handlers"
	"github.com/dro14/yordamchi/recover"
	"github.com/dro14/yordamchi/utils"
)

func main() {
	utils.SetConfigs()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go utils.Main(sigChan)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	recover.Init()
	utils.SendInfoMessage("@yordamchi_ai_bot restarted")
	err := handlers.New().Run(port)
	if err != nil {
		log.Fatal("can't run server:", err)
	}
}
