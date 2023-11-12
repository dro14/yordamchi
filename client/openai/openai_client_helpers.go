package openai

import (
	"context"
	"fmt"
	"github.com/dro14/yordamchi/lib/types"
	"strings"

	"github.com/dro14/yordamchi/lib/models"
	"github.com/dro14/yordamchi/redis"
)

func maxTokens(ctx context.Context) int {

	if ctx.Value("model") == models.GPT3 {
		return 4096
	}

	availableTokens := redis.GPT4Tokens(ctx)
	if availableTokens < 4096 {
		return availableTokens
	} else {
		return 4096
	}
}

func id(ctx context.Context) string {
	return fmt.Sprintf("%d", ctx.Value("user_id").(int64))
}

func vision(messages []types.Message) []types.Message {
	n := len(messages)
	URL, text, _ := strings.Cut(messages[n-1].Content.(string), "\n\n\n")
	var content []types.Content
	if len(text) > 0 {
		content = append(content, types.Content{Type: "text", Text: text})
	}
	content = append(content, types.Content{Type: "image_url", ImageURL: types.ImageURL{URL: URL}})
	messages[n-1] = types.Message{Role: "user", Content: content}
	return messages
}
