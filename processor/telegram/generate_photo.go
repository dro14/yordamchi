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
	prompt, err := translator.Translate("auto", "en", strings.TrimSpace(prompt))

	imageURL := url.QueryEscape(openai.Generations(ctx, prompt))

	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto?chat_id=%d&photo=%s", os.Getenv("BOT_TOKEN"), userID, imageURL))
	if err != nil {
		log.Printf("can't send photo: %v\n", err)
		bts, _ := io.ReadAll(resp.Body)
		log.Printf("body: %s", string(bts))
		return
	}
	_ = resp.Body.Close()

	redis.Decrement(ctx, 0)
}
