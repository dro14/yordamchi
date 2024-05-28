package processor

import (
	"context"
	"log"
	"strings"

	"github.com/dro14/yordamchi/clients/openai/models"
)

var keywords = []string{
	"tarjim", "rus", "ingliz", "ingiliz", "turk", "koreys",
	"перевод", "русск", "английск", "турецк", "корейск",
	"perevod", "russk", "angliysk", "turetsk", "koreys",
	"translat", "russian", "english", "turkish", "korean",
}

func (p *Processor) needTranslation(ctx context.Context, prompt string, userID int64) bool {
	if model(ctx) != models.GPT3 || lang(ctx) != "uz" {
		return false
	}
	lowered := strings.ToLower(prompt)
	for _, keyword := range keywords {
		if strings.Contains(lowered, keyword) {
			log.Printf("user %d asks for translation: %q", userID, prompt)
			return false
		}
	}
	return true
}
