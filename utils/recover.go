package utils

import (
	"log"
	"os"
	"runtime/debug"
)

func LogShutdown(sigChan chan os.Signal) {
	sig := <-sigChan
	log.Printf("Received %v, initiating shutdown...", sig)
	SendLogFile("gin.log")
	SendLogFile("yordamchi.log")
}

func RecoverIfPanic() {
	if r := recover(); r != nil {
		log.Printf("%s\n%s", r, debug.Stack())
		SendLogFile("gin.log")
		SendLogFile("yordamchi.log")
	}
}
