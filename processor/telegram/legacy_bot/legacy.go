package legacy_bot

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gotd/td/tg"
)

var bot *tgbotapi.BotAPI

func Init() {

	token, ok := os.LookupEnv("LEGACY_BOT_TOKEN")
	if !ok {
		log.Fatalf("legacy bot token is not specified")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("can't initialize legacy bot: %v", err)
	}
}

func Reply(c *gin.Context) {

	update := &tgbotapi.Update{}
	if err := c.ShouldBindJSON(update); err != nil {
		log.Printf("can't bind json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if 

	tgbotapi.NewMessage(update.)

	bot.Request()

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
