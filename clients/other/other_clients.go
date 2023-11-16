package other

import (
	"log"
	"os"
)

type APIs struct {
	subscriptionKey string
	visionURL       string
	translateURL    string
	searchURL       string
}

func New() *APIs {
	subscriptionKey, ok := os.LookupEnv("SUBSCRIPTION_KEY")
	if !ok {
		log.Fatal("subscription key is not specified")
	}

	return &APIs{
		subscriptionKey: subscriptionKey,
		visionURL:       "https://yordamchi.cognitiveservices.azure.com/vision/v3.1/ocr",
		translateURL:    "https://translate.google.com/m?sl=%s&tl=%s&q=%s",
		searchURL:       "https://google.victoriousriver-fffd2d70.westeurope.azurecontainerapps.io/search",
	}
}
