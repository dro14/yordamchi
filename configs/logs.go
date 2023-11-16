package configs

import (
	"log"
	"os"
	"runtime/debug"
)

func Main(sigChan chan os.Signal) {
	sig := <-sigChan
	log.Printf("Received %v, initiating shutdown...", sig)
	SendFile("gin.log")
	SendFile("yordamchi.log")
}

func RecoverIfPanic() {
	if r := recover(); r != nil {
		log.Printf("%s\n%s", r, debug.Stack())
		SendFile("gin.log")
		SendFile("yordamchi.log")
	}
}
