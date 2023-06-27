package openai

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/dro14/yordamchi/client/bobdev"
	"github.com/dro14/yordamchi/client/openai"
	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/text"
)

func Process(ctx context.Context, messages []types.Message, stats *types.Stats, channel chan<- string) {

	tokenLimit := 0
	if ctx.Value("model") == "gpt-4" {
		tokenLimit = 8191
	} else {
		tokenLimit = 4096
	}

	maxTokens := tokens(ctx, messages)
	retryDelay := 10 * constants.RetryDelay
	var errMsg string
Retry:
	stats.Attempts++
	response, err := openai.Completion(ctx, messages, maxTokens, channel)
	if err != nil {
		errMsg = err.Error()

		switch {
		case strings.HasPrefix(errMsg, e.InvalidRequest):
			errMsg = strings.TrimPrefix(errMsg, e.InvalidRequest)
			if strings.HasPrefix(errMsg, e.ContextLengthExceeded) {
				errMsg = strings.TrimPrefix(errMsg, e.ContextLengthExceeded)
				errMsg, _, _ = strings.Cut(errMsg, " tokens")
				totalTokens, _ := strconv.Atoi(errMsg)
				diff := totalTokens - tokenLimit
				log.Printf("max tokens %d was deacreased by %d", maxTokens, diff)
				maxTokens -= diff
			} else if len(messages) > 2 {
				messages = messages[2:]
				maxTokens = tokens(ctx, messages)
			} else {
				channel <- text.TooLong[lang(ctx)]
				return
			}
		case strings.HasPrefix(errMsg, e.EmptyCompletion):
			channel <- text.RequestFailed[lang(ctx)]
			return
		case strings.HasPrefix(errMsg, e.StreamError):
			channel <- text.Error[lang(ctx)]
		case strings.HasPrefix(errMsg, e.ServiceUnavailable),
			strings.HasPrefix(errMsg, e.InternalServerError):
			functions.Sleep(&retryDelay)
		}

		if stats.Attempts < constants.RetryAttempts {
			goto Retry
		} else {
			log.Printf("%q failed after %d attempts", errMsg, stats.Attempts)
			channel <- text.RequestFailed[lang(ctx)]
			return
		}
	} else if stats.Attempts > 1 {
		log.Printf("%q was handled after %d attempts", errMsg, stats.Attempts)
	}

	stats.FinishReason = response.Choices[0].FinishReason
	stats.PromptTokens = tokenLimit - maxTokens
	stats.PromptLength = length(messages)

	completions := []types.Message{response.Choices[0].Message}
	stats.CompletionTokens = bobdev.Tokens(ctx, completions) - 8
	stats.CompletionLength = len(completions[0].Content)
}
