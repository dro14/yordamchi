package functions

import (
	"log"
	"time"
)

func LanguageCode(lang string) string {
	if lang == "" {
		lang = "uz"
	} else if lang != "uz" && lang != "ru" {
		lang = "en"
	}
	return lang
}

func Sleep(retryDelay *time.Duration) {
	if *retryDelay > 0 {
		log.Printf("retrying request after %v", *retryDelay)
		time.Sleep(*retryDelay)
		*retryDelay *= 2
	}
}
