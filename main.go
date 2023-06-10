package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	redisCache "github.com/dro14/yordamchi/cache/redis"
	tgProcessor "github.com/dro14/yordamchi/processor/telegram"
	"github.com/dro14/yordamchi/processor/telegram/info"
	"github.com/dro14/yordamchi/processor/telegram/legacy"
	"github.com/gotd/contrib/redis"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	time.Local, _ = time.LoadLocation("Asia/Tashkent")

	dispatcher := tg.NewUpdateDispatcher()
	cache := redisCache.New()
	session := redis.NewSessionStorage(cache.Redis, "main_bot_session")

	if err := telegram.BotFromEnvironment(
		context.Background(),
		telegram.Options{
			UpdateHandler:  dispatcher,
			SessionStorage: session,
		},
		func(ctx context.Context, client *telegram.Client) error {

			processor := tgProcessor.New(client.API(), cache)

			dispatcher.OnNewMessage(processor.ProcessMessage)
			dispatcher.OnBotCallbackQuery(processor.ProcessCallbackQuery)
			dispatcher.OnBotStopped(processor.ProcessBotStopped)

			go processor.Recover()
			go legacy.ConnectLegacyBot(cache.Redis)
			go info.ConnectInfoBot(cache.Redis)
			go listen()

			log.Printf("main bot is connected")
			return nil
		},
		telegram.RunUntilCanceled,
	); err != nil {
		log.Fatalf("can't connect client: %v", err)
	}
}

func listen() {

	port, ok := os.LookupEnv("PORT")
	if !ok {
		return
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("can't start server: %v", err)
	}
}
