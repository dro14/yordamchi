package info

import (
	"context"
	"log"
	"os"
	"time"

	redisClient "github.com/go-redis/redis/v8"
	"github.com/gotd/contrib/redis"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

var infoBot *tg.Client

func ConnectInfoBot(cacheClient *redisClient.Client) {

	token, ok := os.LookupEnv("INFO_BOT_TOKEN")
	if !ok {
		log.Fatalf("info bot token is not specified")
	}

	session := redis.NewSessionStorage(cacheClient, "info_bot_session")

	client, err := telegram.ClientFromEnvironment(
		telegram.Options{
			SessionStorage: session,
		},
	)
	if err != nil {
		log.Fatalf("can't create client: %v", err)
	}

	if err = client.Run(context.Background(), func(ctx context.Context) error {

		status, err := client.Auth().Status(ctx)
		if err != nil {
			log.Fatalf("can't get authorization status: %v", err)
		}

		if !status.Authorized {
			_, err = client.Auth().Bot(ctx, token)
			if err != nil {
				log.Fatalf("can't authorize bot: %v", err)
			}
		}

		infoBot = client.API()
		SendInfoMessage(ctx, "@yordamchi_ai_bot restarted")

		log.Printf("info bot is connected")
		return telegram.RunUntilCanceled(ctx, client)
	}); err != nil {
		log.Fatalf("can't connect info bot: %v", err)
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
