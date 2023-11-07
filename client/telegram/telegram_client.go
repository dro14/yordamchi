package telegram

import (
	"context"
	"log"
	"os"

	"github.com/dro14/yordamchi/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	cache "github.com/gotd/contrib/redis"
	"github.com/gotd/td/telegram"
)

func Init() {

	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatalf("bot token is not specified")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("can't initialize bot: %v", err)
	}

	sessionStorage := cache.NewSessionStorage(redis.Client, "main_bot_session")
	done := make(chan bool)

	go func() {
		if err = telegram.BotFromEnvironment(
			context.Background(),
			telegram.Options{SessionStorage: sessionStorage},
			func(ctx context.Context, client *telegram.Client) error {
				api = client.API()
				done <- true
				return nil
			},
			telegram.RunUntilCanceled,
		); err != nil {
			log.Fatalf("can't connect client: %v", err)
		}
	}()

	<-done
}
