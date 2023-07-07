package legacy_bot

import (
	"log"
	"os"

	"github.com/dro14/yordamchi/text"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true})

	if update.Message == nil {
		return
	}

	lang := update.Message.From.LanguageCode
	if lang == "" {
		lang = "uz"
	} else if lang != "uz" && lang != "ru" {
		lang = "en"
	}

	config := tgbotapi.NewMessage(update.Message.From.ID, text.LegacyMessage[lang])
	_, err := bot.Send(config)
	if err != nil {
		log.Printf("can't send legacy message: %v", err)
	}
}
