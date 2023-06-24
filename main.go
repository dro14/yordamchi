package main

import (
	"context"
	"log"
	"time"

	"github.com/dro14/yordamchi/payme"
	tgProcessor "github.com/dro14/yordamchi/processor/telegram"
	"github.com/dro14/yordamchi/processor/telegram/info"
	"github.com/dro14/yordamchi/processor/telegram/legacy"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {

	time.Local, _ = time.LoadLocation("Asia/Tashkent")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dispatcher := tg.NewUpdateDispatcher()
	if err := telegram.BotFromEnvironment(
		context.Background(),
		telegram.Options{UpdateHandler: dispatcher},
		func(ctx context.Context, client *telegram.Client) error {

			processor := tgProcessor.New(client.API())

			dispatcher.OnNewMessage(processor.ProcessMessage)
			dispatcher.OnBotCallbackQuery(processor.ProcessCallbackQuery)
			dispatcher.OnBotStopped(processor.ProcessBotStopped)

			go processor.Recover()
			go legacy.Run()
			go info.Run()
			go payme.Run()

			log.Printf("main bot is connected")
			return nil
		},
		telegram.RunUntilCanceled,
	); err != nil {
		log.Fatalf("can't connect client: %v", err)
	}
}
