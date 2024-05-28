package processor

import (
	"context"
	"log"
	"strings"

	"github.com/dro14/yordamchi/clients/openai/models"
)

var keywords = []string{
	"tarjim", "o'zbek til", "ozbek til", "rus til", "ingliz til", "ingiliz til", "turk til", "koreys til",
	"o'zbekcha", "ozbekcha", "ruscha", "inglizcha", "ingilizcha", "turkcha", "koreyscha",
	"перевод", "узбекск", "русск", "английск", "турецк", "корейск",
	"perevod", "uzbeksk", "russk", "angliysk", "turetsk", "koreysk",
	"translat", "uzbek", "russian", "english", "turkish", "korean",
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
