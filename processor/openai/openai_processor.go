package openai

import (
	"context"
	"github.com/dro14/yordamchi/lib/functions"
	"log"
	"strconv"
	"strings"

	"github.com/dro14/yordamchi/client/bobdev"
	"github.com/dro14/yordamchi/client/openai"
	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/processor/telegram/text"
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
	retryDelay := constants.RetryDelay
	var errMsg string
	var found bool
Retry:
	stats.Attempts++
	response, err := p.aiClient.Completion(ctx, messages, maxTokens, channel)
	if err != nil {
		errMsg = err.Error()

		if errMsg, found = strings.CutPrefix(errMsg, e.InvalidRequest); found {
			if errMsg, found = strings.CutPrefix(errMsg, e.ContextLengthExceeded); found {
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
		} else if strings.HasPrefix(errMsg, e.EmptyCompletion) {
			channel <- text.RequestFailed[lang(ctx)]
			return
		} else if strings.HasPrefix(errMsg, e.StreamError) {
			channel <- text.Error[lang(ctx)]
		} else if strings.HasPrefix(errMsg, e.ServiceUnavailable) {
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
