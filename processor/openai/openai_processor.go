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

type Processor struct {
	aiClient *openai.Client
}

func New() *Processor {
	return &Processor{
		aiClient: openai.New(),
	}
}

func (p *Processor) Process(ctx context.Context, messages []types.Message, stats *types.Stats, channel chan<- string) {

	maxTokens := 4096 - bobdev.Tokens("gpt-3.5-turbo", messages)
	retryDelay := 10 * constants.RetryDelay
	var errMsg string
Retry:
	stats.Attempts++
	response, err := p.aiClient.Completion(ctx, messages, maxTokens, channel)
	if err != nil {
		errMsg = err.Error()

		switch {
		case strings.HasPrefix(errMsg, e.InvalidRequest):
			errMsg = strings.TrimPrefix(errMsg, e.InvalidRequest)
			if strings.HasPrefix(errMsg, e.ContextLengthExceeded) {
				errMsg = strings.TrimPrefix(errMsg, e.ContextLengthExceeded)
				errMsg, _, _ = strings.Cut(errMsg, " tokens")
				tokens, _ := strconv.Atoi(errMsg)
				diff := tokens - 4096
				log.Printf("max tokens %d was deacreased by %d", maxTokens, diff)
				maxTokens -= diff
			} else if len(messages) > 2 {
				messages = messages[2:]
				maxTokens = 4096 - bobdev.Tokens("gpt-3.5-turbo", messages)
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
	stats.PromptTokens = 4096 - maxTokens
	stats.PromptLength = length(messages)

	completions := []types.Message{response.Choices[0].Message}
	stats.CompletionTokens = bobdev.Tokens("gpt-3.5-turbo", messages) - 8
	stats.CompletionLength = len(completions[0].Content)
}
