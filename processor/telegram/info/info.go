package info

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gotd/td/tg"
)

var bot *tgbotapi.BotAPI

func Run() {

	token, ok := os.LookupEnv("INFO_BOT_TOKEN")
	if !ok {
		log.Fatalf("info bot token is not specified")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("can't initialize info bot: %v", err)
	}
}

func SendInfoMessage(ctx context.Context, text string) {

	request := &tg.MessagesSendMessageRequest{
		Peer:      &tg.InputPeerUser{UserID: 1331278972},
		Message:   text,
		RandomID:  time.Now().UnixNano(),
		NoWebpage: true,
	}

	_, err := infoBot.MessagesSendMessage(ctx, request)
	if err != nil {
		log.Printf("can't send info message: %v", err)
	}

	request.Peer = &tg.InputPeerUser{UserID: 835282186}
	request.RandomID = time.Now().UnixNano()

	_, err = infoBot.MessagesSendMessage(ctx, request)
	if err != nil {
		log.Printf("can't send info message: %v", err)
	}
}
