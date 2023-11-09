package telegram

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dro14/yordamchi/client/openai"
	"github.com/dro14/yordamchi/client/translator"
	"github.com/dro14/yordamchi/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GeneratePhoto(ctx context.Context, message *tgbotapi.Message) {
	userID := ctx.Value("user_id").(int64)

	prompt := strings.ReplaceAll(message.Text, "#image", "")
	prompt = strings.TrimSpace(prompt)
	prompt, _ = translator.Translate("auto", "en", prompt)

	imageURL := openai.Generations(ctx, prompt)
	imageURL = url.QueryEscape(imageURL)
	botToken := os.Getenv("BOT_TOKEN")
	URL := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto?chat_id=%d&photo=%s", botToken, userID, imageURL)

	resp, err := http.Get(URL)
	if err != nil {
		bts, _ := io.ReadAll(resp.Body)
		log.Printf("can't send photo: %v\nbody: %s", err, string(bts))
		return
	}
	_ = resp.Body.Close()

	redis.Decrement(ctx, 0)
}
