package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dro14/yordamchi/configs"
	"github.com/dro14/yordamchi/handlers"
	"github.com/dro14/yordamchi/recover"
)

func main() {
	configs.Init()
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer configs.Main(sigChan)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	recover.Init()
	configs.SendMessage("@yordamchi_ai_bot restarted")
	err := handlers.New().Run(port)
	if err != nil {
		log.Fatal("can't run server:", err)
	}
}
