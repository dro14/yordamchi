package utils

import (
	"log"
	"os"
	"runtime/debug"
)

func LogShutdown(sigChan chan os.Signal) {
	sig := <-sigChan
	log.Printf("Received %v, initiating shutdown...", sig)
	SendLogFile("yordamchi.log")
	SendLogFile("gin.log")
}

func RecoverIfPanic() {
	if r := recover(); r != nil {
		log.Printf("%s\n%s", r, debug.Stack())
		SendLogFile("yordamchi.log")
		SendLogFile("gin.log")
	}
}
