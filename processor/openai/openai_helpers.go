package openai

import (
	"context"
	"fmt"
	"github.com/dro14/yordamchi/client/bobdev"
	"github.com/dro14/yordamchi/redis"

	"github.com/dro14/yordamchi/lib/types"
)

func length(messages []types.Message) int {
	var promptLength int
	for i := range messages {
		promptLength += len(fmt.Sprintf("role: %s\ncontent: %s", messages[i].Role, messages[i].Content))
	}
	return promptLength
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

func tokens(ctx context.Context, messages []types.Message) int {

	if ctx.Value("model") == "gpt-3.5-turbo" {
		return 4096 - bobdev.Tokens(ctx, messages)
	}

	maxTokens := 4096 - bobdev.Tokens(ctx, messages)
	availableTokens := redis.GPT4Tokens(ctx)
	if availableTokens < maxTokens {
		return availableTokens
	} else {
		return maxTokens
	}
}
