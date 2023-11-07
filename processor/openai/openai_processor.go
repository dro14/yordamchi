package openai

import (
	"context"
	"log"
	"strings"

	"github.com/dro14/yordamchi/client/bobdev"
	"github.com/dro14/yordamchi/client/openai"
	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/text"
)

func ProcessWithStream(ctx context.Context, messages []types.Message, stats *types.Stats, channel chan<- string) {

	retryDelay := 10 * constants.RetryDelay
	var errMsg string
Retry:
	stats.Attempts++
	response, err := openai.CompletionWithStream(ctx, messages, channel)
	if err != nil {
		errMsg = err.Error()

		switch {
		case strings.HasPrefix(errMsg, e.InvalidRequest):
			errMsg = strings.TrimPrefix(errMsg, e.InvalidRequest)
			if strings.HasPrefix(errMsg, e.ContextLengthExceededGPT3) ||
				strings.HasPrefix(errMsg, e.ContextLengthExceededGPT4) {
				errMsg = strings.TrimPrefix(errMsg, e.ContextLengthExceededGPT3)
				errMsg = strings.TrimPrefix(errMsg, e.ContextLengthExceededGPT4)
				errMsg, _, _ = strings.Cut(errMsg, " tokens")
			} else if len(messages) > 2 {
				messages = messages[2:]
			} else {
				channel <- text.TooLong[lang(ctx)]
				return
			}
			goto Retry
		case strings.HasPrefix(errMsg, e.StreamError):
			channel <- text.Error[lang(ctx)]
			goto Retry
		case strings.HasPrefix(errMsg, e.BadGateway),
			strings.HasPrefix(errMsg, e.ServiceUnavailable),
			strings.Contains(errMsg, e.ContextDeadlineExceeded):
			goto Retry
		}

		if stats.Attempts < constants.RetryAttempts {
			functions.Sleep(&retryDelay)
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
	stats.PromptTokens = bobdev.Tokens(ctx, messages)
	stats.PromptLength = length(messages)

	completions := []types.Message{response.Choices[0].Message}
	stats.CompletionTokens = bobdev.Tokens(ctx, completions) - 8
	stats.CompletionLength = len(completions[0].Content)
}

func Process(ctx context.Context, messages []types.Message, stats *types.Stats) (string, error) {

	retryDelay := 10 * constants.RetryDelay
	var errMsg string
Retry:
	stats.Attempts++
	response, err := openai.Completion(ctx, messages)
	if err != nil {
		errMsg = err.Error()

		switch {
		case strings.HasPrefix(errMsg, e.InvalidRequest):
			errMsg = strings.TrimPrefix(errMsg, e.InvalidRequest)
			if strings.HasPrefix(errMsg, e.ContextLengthExceededGPT3) ||
				strings.HasPrefix(errMsg, e.ContextLengthExceededGPT4) {
				errMsg = strings.TrimPrefix(errMsg, e.ContextLengthExceededGPT3)
				errMsg = strings.TrimPrefix(errMsg, e.ContextLengthExceededGPT4)
				errMsg, _, _ = strings.Cut(errMsg, " tokens")
			} else if len(messages) > 2 {
				messages = messages[2:]
			} else {
				log.Printf("%s", errMsg)
				return text.TooLong[lang(ctx)], err
			}
			goto Retry
		case strings.HasPrefix(errMsg, e.BadGateway),
			strings.HasPrefix(errMsg, e.ServiceUnavailable),
			strings.Contains(errMsg, e.ContextDeadlineExceeded):
			goto Retry
		}

		if stats.Attempts < constants.RetryAttempts {
			functions.Sleep(&retryDelay)
			goto Retry
		} else {
			log.Printf("%q failed after %d attempts", errMsg, stats.Attempts)
			return text.RequestFailed[lang(ctx)], err
		}
	} else if stats.Attempts > 1 {
		log.Printf("%q was handled after %d attempts", errMsg, stats.Attempts)
	}

	stats.FinishReason = response.Choices[0].FinishReason
	stats.PromptTokens = bobdev.Tokens(ctx, messages)
	stats.PromptLength = length(messages)

	completions := []types.Message{response.Choices[0].Message}
	stats.CompletionTokens = bobdev.Tokens(ctx, completions) - 8
	stats.CompletionLength = len(completions[0].Content)
	return completions[0].Content, nil
}
