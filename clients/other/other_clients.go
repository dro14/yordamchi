package other

import (
	"log"
	"os"
)

type APIs struct {
	subscriptionKey string
	ocrURL          string
	translateURL    string
}

func New() *APIs {
	subscriptionKey, ok := os.LookupEnv("SUBSCRIPTION_KEY")
	if !ok {
		log.Fatal("subscription key is not specified")
	}

	return &APIs{
		subscriptionKey: subscriptionKey,
		ocrURL:          "https://yordamchi.cognitiveservices.azure.com/computervision/imageanalysis:analyze?api-version=2022-10-12-preview&features=read",
		translateURL:    "https://translate.google.com/m?sl=%s&tl=%s&q=%s",
	}
}
