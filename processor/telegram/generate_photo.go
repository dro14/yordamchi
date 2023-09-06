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
	"github.com/dro14/yordamchi/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GeneratePhoto(ctx context.Context, message *tgbotapi.Message) {

	userID := ctx.Value("user_id").(int64)

	prompt := strings.ReplaceAll(message.Text, "#image", "")
	prompt = strings.TrimSpace(prompt)

	imageURL := url.QueryEscape(openai.Generations(ctx, prompt))
	fmt.Println(imageURL)

	resp, err := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto?chat_id=%d&photo=%s", os.Getenv("BOT_TOKEN"), userID, imageURL))
	if err != nil {
		log.Printf("can't send photo: %v", err)
		return
	}
	defer resp.Body.Close()

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("can't read response: %v", err)
		return
	}

	fmt.Println(string(bts))
	redis.Decrement(ctx, 0)
}
