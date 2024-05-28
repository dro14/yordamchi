package processor

import (
	"context"
	"encoding/json"
	"log"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/clients/openai/types"
)

var system = `DOES THE USER ASK FOR TRANSLATION?
RESPOND IN JSON FORMAT:

{
    "answer": <boolean>
}`

func (p *Processor) needTranslation(ctx context.Context, prompt string, userID int64) bool {

	if model(ctx) != models.GPT3 || lang(ctx) != "uz" {
		return false
	}

	messages := []types.Message{
		{Role: "system", Content: system},
		{Role: "user", Content: prompt},
	}

	ctx = context.WithValue(ctx, "stream", false)
	ctx = context.WithValue(ctx, "json_mode", true)

	response, err := p.openai.Completions(ctx, messages, nil, "", make(chan string, 1))
	if err != nil {
		log.Printf("user %d: can't check if the user asks for translation: %v", userID, err)
		return true
	}

	bts := []byte(response.Choices[0].Message.Content.(string))
	var answer map[string]bool
	err = json.Unmarshal(bts, &answer)
	if err != nil {
		log.Printf("user %d: can't decode response: %v\nbody: %s", userID, err, bts)
		return true
	}

	if answer["answer"] {
		log.Printf("user %d asks for translation: %q", userID, prompt)
		return false
	} else {
		return true
	}
}
