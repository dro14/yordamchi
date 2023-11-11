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
)

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
}

func Read(ctx context.Context, photoURL, caption string) string {
	var buffer bytes.Buffer
	request := map[string]string{"url": photoURL}
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Printf("can't encode request: %v", err)
		return caption
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://yordamchi.cognitiveservices.azure.com/vision/v3.1/ocr", &buffer)
	if err != nil {
		log.Printf("can't create request: %v", err)
		return caption
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", constants.SubscriptionKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("can't do request: %v", err)
		return caption
	}
	defer func() { _ = resp.Body.Close() }()

	response := &Response{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		log.Printf("can't decode response: %v", err)
		return caption
	}

	var builder strings.Builder
	if len(caption) > 0 {
		builder.WriteString(caption + ":\n\n")
	}

	for _, region := range response.Regions {
		for _, line := range region.Lines {
			for _, word := range line.Words {
				builder.WriteString(word.Text + " ")
			}
			builder.WriteString("\n")
		}
		builder.WriteString("\n\n")
	}

	return builder.String()
}
