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

	recover.Start()
	utils.SendInfoMessage(-1001924963699, 35)
	err := handlers.New().Run(port)
	if err != nil {
		log.Fatal("can't run server:", err)
	}
}
