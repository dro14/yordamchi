package redis

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/clients/openai/types"
)

func (r *Redis) ConversationHistory(ctx context.Context, prompt string) []types.Message {
	var message types.Message
	if ctx.Value("model") == models.GPT4V {
		URL, text, _ := strings.Cut(prompt, "\n\n\n")
		content := []types.Content{{Type: "image_url", ImageURL: &types.ImageURL{URL: URL}}}
		if len(text) > 0 {
			content = append(content, types.Content{Type: "text", Text: text})
		}
		message = types.Message{Role: "user", Content: content}
	} else {
		message = types.Message{Role: "user", Content: prompt}
	}

	var messages []types.Message
	jsonData, err := r.client.Get(ctx, "context:"+id(ctx)).Bytes()
	if err != nil {
		return append(messages, message)
	}

	err = json.Unmarshal(jsonData, &messages)
	if err != nil {
		log.Printf("can't decode %q: %s", "context:"+id(ctx), err)
		return append(messages, message)
	}
	return append(messages, message)
}

func (r *Redis) StoreHistory(ctx context.Context, messages []types.Message, completion string) {
	messages = append(messages, types.Message{Role: "assistant", Content: completion})
	messages = messages[len(messages)-2:]

	for _, message := range messages {
		content, ok := message.Content.([]types.Content)
		if ok {
			if len(content) == 2 {
				message.Content = content[1].Text
			} else {
				message.Content = ""
			}
		}
	}

	jsonData, err := json.Marshal(messages)
	if err != nil {
		log.Printf("can't encode %q: %s", "context:"+id(ctx), err)
		return
	}
	r.client.Set(ctx, "context:"+id(ctx), jsonData, 72*time.Hour)
}

func (r *Redis) DeleteHistory(ctx context.Context) {
	r.client.Del(ctx, "context:"+id(ctx))
}
