package openai

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/dro14/yordamchi/processor/text"
	"github.com/dro14/yordamchi/storage/postgres"
	"github.com/dro14/yordamchi/utils"
)

func (o *OpenAI) ProcessCompletions(ctx context.Context, prompt string, stats *postgres.Stats, channel chan<- string) {
	defer close(channel)
	defer utils.RecoverIfPanic()
	messages := o.redis.ConversationHistory(ctx, prompt)
	retryDelay := 10 * utils.RetryDelay
	var errMsg string
Retry:
	stats.Attempts++
	response, err := o.Completions(ctx, messages, channel)
	if err != nil {
		errMsg = err.Error()
		if strings.Contains(errMsg, "stream error") {
			channel <- text.Error[lang(ctx)]
			retryDelay = 0
		} else if strings.Contains(errMsg, "context deadline exceeded") {
			retryDelay = 0
		}
		if stats.Attempts < utils.RetryAttempts {
			utils.Sleep(&retryDelay)
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
	stats.PromptTokens = o.countTokens(messages)
	stats.PromptLength = length(messages)

	completion := response.Choices[0].Message.Content.(string)
	stats.CompletionTokens = o.countTokens(completion)
	stats.CompletionLength = len(completion)

	o.redis.StoreHistory(ctx, messages, completion)
	time.Sleep(utils.RequestInterval)
	channel <- completion
}

func (o *OpenAI) ProcessGenerations(ctx context.Context, prompt string) (string, error) {
	retryDelay := 10 * utils.RetryDelay
	var errMsg string
	var attempts int
Retry:
	attempts++
	response, err := o.Generations(ctx, prompt)
	if err != nil {
		errMsg = err.Error()
		if strings.Contains(errMsg, "context deadline exceeded") {
			retryDelay = 0
		}
		if attempts < utils.RetryAttempts {
			utils.Sleep(&retryDelay)
			goto Retry
		} else {
			log.Printf("%q failed after %d attempts", errMsg, attempts)
			return text.RequestFailed[lang(ctx)], err
		}
	} else if attempts > 1 {
		log.Printf("%q was handled after %d attempts", errMsg, attempts)
	}

	return response.Data[0].URL, nil
}
