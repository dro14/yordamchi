package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Regions []struct {
		Lines []struct {
			Words []struct {
				Text string `json:"text"`
			} `json:"words"`
		} `json:"lines"`
	} `json:"regions"`
}

func Init() {

	subscriptionKey, ok := os.LookupEnv("SUBSCRIPTION_KEY")
	if !ok {
		log.Fatalf("subscription key is not specified")
	}
	constants.SubscriptionKey = subscriptionKey

	var err error
	bot, err = tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatalf("can't initialize bot: %v", err)
	}
}

func Analyze(ctx context.Context, message *tgbotapi.Message) string {

	photo := message.Photo[len(message.Photo)-1]
	fileURL, err := bot.GetFileDirectURL(photo.FileID)
	if err != nil {
		log.Printf("can't get file url: %v", err)
		return message.Caption
	}

	request := &Request{fileURL}
	var buffer bytes.Buffer
	err = json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Printf("can't encode request: %v", err)
		return message.Caption
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://yordamchi.cognitiveservices.azure.com/vision/v3.1/ocr", &buffer)
	if err != nil {
		log.Printf("can't create request: %v", err)
		return message.Caption
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", constants.SubscriptionKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("can't do request: %v", err)
		return message.Caption
	}
	defer func() { _ = resp.Body.Close() }()

	response := &Response{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		log.Printf("can't decode response: %v", err)
		return message.Caption
	}

	var builder strings.Builder
	if len(message.Caption) > 0 {
		builder.WriteString(message.Caption + ":\n\n")
	}

	for _, region := range response.Regions {
		for _, line := range region.Lines {
			for _, word := range line.Words {
				builder.WriteString(word.Text + " ")
			}
			builder.WriteString("\n")
		}
	}

	return builder.String()
}
