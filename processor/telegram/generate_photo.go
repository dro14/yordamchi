package telegram

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

	imageURL := openai.Generations(ctx, prompt)

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto?chat_id=%d&photo=%s", os.Getenv("BOT_TOKEN"), userID, imageURL)
	_, err := http.Get(url)
	if err != nil {
		log.Printf("can't send photo: %v", err)
	}

	redis.Decrement(ctx, 0)
}
