package openai

import (
	"context"
	"fmt"

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
