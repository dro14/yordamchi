package main

import (
	"context"
	"log"
	"time"

	"github.com/dro14/yordamchi/payme"
	"github.com/dro14/yordamchi/postgres"
	tgProcessor "github.com/dro14/yordamchi/processor/telegram"
	"github.com/dro14/yordamchi/processor/telegram/info"
	"github.com/dro14/yordamchi/processor/telegram/legacy"
	"github.com/dro14/yordamchi/redis"
	cache "github.com/gotd/contrib/redis"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	time.Local, _ = time.LoadLocation("Asia/Tashkent")
	redis.Init()
	postgres.Init()

	dispatcher := tg.NewUpdateDispatcher()
	if err := telegram.BotFromEnvironment(
		context.Background(),
		telegram.Options{
			UpdateHandler:  dispatcher,
			SessionStorage: cache.NewSessionStorage(redis.Client, "main_bot_session"),
		},
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
