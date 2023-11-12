package recover

import (
	"log"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/dro14/yordamchi/processor/telegram/info_bot"
)

func Recover(sigChan chan os.Signal) {
	sig := <-sigChan
	if r := recover(); r != nil {
		log.Printf("%s\n%s", r, debug.Stack())
	}
	log.Printf("Received %q, initiating shutdown...", sig)
	info_bot.SendFile("yordamchi.log")
	os.Exit(0)
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
