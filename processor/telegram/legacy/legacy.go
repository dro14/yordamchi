package legacy

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/text"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

var legacyBot *tg.Client

func Run() {

	token, ok := os.LookupEnv("LEGACY_BOT_TOKEN")
	if !ok {
		log.Fatalf("legacy bot token is not specified")
	}

	dispatcher := tg.NewUpdateDispatcher()
	client, err := telegram.ClientFromEnvironment(telegram.Options{UpdateHandler: dispatcher})
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

		legacyBot = client.API()
		dispatcher.OnNewMessage(SendLegacyMessage)

		log.Printf("legacy bot is connected")
		return telegram.RunUntilCanceled(ctx, client)
	}); err != nil {
		log.Fatalf("can't connect legacy bot: %v", err)
	}
}

func SendLegacyMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) error {

	message, ok := update.Message.(*tg.Message)
	if !ok || message.Out {
		return nil
	}

	peerUser, ok := message.PeerID.(*tg.PeerUser)
	if !ok {
		return nil
	}

	userID := peerUser.UserID
	user := entities.Users[userID]
	langCode := functions.LanguageCode(user.LangCode)

	request := &tg.MessagesSendMessageRequest{
		Peer:      &tg.InputPeerUser{UserID: userID},
		Message:   text.LegacyMessage[langCode],
		RandomID:  time.Now().UnixNano(),
		NoWebpage: true,
	}

	_, err := legacyBot.MessagesSendMessage(ctx, request)
	if err != nil {
		log.Printf("can't send legacy message: %v", err)
	}

	return nil
}
