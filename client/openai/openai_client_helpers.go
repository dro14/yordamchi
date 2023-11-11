package openai

import (
	"context"

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
