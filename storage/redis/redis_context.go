package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/clients/openai/types"
	"github.com/dro14/yordamchi/utils"
)

const template = "You are a friendly chatbot in Telegram called Yordamchi, based on %s model architecture."

func (r *Redis) ConversationHistory(ctx context.Context, prompt string) (output context.Context, messages []types.Message) {
	systemMessage := types.Message{Role: "system"}
	if ctx.Value("model") == models.GPT3 {
		systemMessage.Content = fmt.Sprintf(template, "GPT-3.5")
	} else {
		systemMessage.Content = fmt.Sprintf(template, "GPT-4")
	}
	userMessage := types.Message{Role: "user", Content: prompt}

	defer func() {
		messages = append([]types.Message{systemMessage}, messages...)
		messages = append(messages, userMessage)
		for i := range messages {
			URL, text, found := strings.Cut(messages[i].Content.(string), utils.Delimiter)
			if !found {
				continue
			}
			var content []types.Content
			if len(text) > 0 {
				content = append(content, types.Content{Type: "text", Text: text})
			}
			content = append(content, types.Content{Type: "image_url", ImageURL: &types.ImageURL{URL: URL}})
			messages[i].Content = content
			output = context.WithValue(ctx, "model", models.GPT4V)
		}
		if output == nil {
			output = ctx
		}
	}()

	jsonData, err := r.client.Get(ctx, "context:"+id(ctx)).Bytes()
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonData, &messages)
	if err != nil {
		log.Printf("can't decode %q: %s", "context:"+id(ctx), err)
	}
	return
}

func (r *Redis) StoreHistory(ctx context.Context, prompt, completion string) {
	messages := []types.Message{
		{Role: "user", Content: prompt},
		{Role: "assistant", Content: completion},
	}

	jsonData, err := json.Marshal(messages)
	if err != nil {
		log.Printf("can't encode %q: %s", "context:"+id(ctx), err)
		return
	}
	var expiration time.Duration
	if strings.Contains(prompt, utils.Delimiter) {
		expiration = 1 * time.Hour
	} else {
		expiration = 72 * time.Hour
	}
	r.client.Set(ctx, "context:"+id(ctx), jsonData, expiration)
}

func (r *Redis) DeleteHistory(ctx context.Context) {
	r.client.Del(ctx, "context:"+id(ctx))
}
