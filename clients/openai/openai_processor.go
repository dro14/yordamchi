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

func (o *OpenAI) ProcessCompletions(ctx context.Context, prompt string, msg *postgres.Message, channel chan<- string) {
	defer close(channel)
	defer utils.RecoverIfPanic()
	ctx, messages := o.redis.ConversationHistory(ctx, prompt)
	retryDelay := 10 * utils.RetryDelay
	var errMsg string
Retry:
	msg.Attempts++
	response, err := o.Completions(ctx, messages, channel)
	if err != nil {
		errMsg = err.Error()
		if strings.Contains(errMsg, "400 Bad Request") {
			channel <- text.BadRequest[lang(ctx)]
			return
		} else if strings.Contains(errMsg, "stream error") {
			channel <- text.Error[lang(ctx)]
			retryDelay = 0
		} else if strings.Contains(errMsg, "context deadline exceeded") {
			retryDelay = 0
		}
		if msg.Attempts < utils.RetryAttempts {
			utils.Sleep(&retryDelay)
			goto Retry
		} else {
			log.Printf("%q failed after %d attempts", errMsg, msg.Attempts)
			channel <- text.RequestFailed[lang(ctx)]
			return
		}
	} else if msg.Attempts > 1 {
		log.Printf("%q was handled after %d attempts", errMsg, msg.Attempts)
	}

	msg.FinishReason = response.Choices[0].FinishReason
	msg.PromptTokens = o.countTokens(messages)
	msg.PromptLength = length(messages)

	completion := response.Choices[0].Message.Content.(string)
	msg.CompletionTokens = o.countTokens(completion)
	msg.CompletionLength = len(completion)

	o.redis.StoreHistory(ctx, prompt, completion)
	time.Sleep(utils.ReqInterval)
	channel <- completion
}

func (o *OpenAI) ProcessGenerations(ctx context.Context, prompt string) (string, string) {
	retryDelay := 10 * utils.RetryDelay
	var errMsg string
	var attempts int
Retry:
	attempts++
	response, err := o.Generations(ctx, prompt)
	if err != nil {
		errMsg = err.Error()
		if strings.Contains(errMsg, "400 Bad Request") {
			return "", text.BadRequest[lang(ctx)]
		} else if strings.Contains(errMsg, "context deadline exceeded") {
			retryDelay = 0
		}
		if attempts < utils.RetryAttempts {
			utils.Sleep(&retryDelay)
			goto Retry
		} else {
			log.Printf("%q failed after %d attempts", errMsg, attempts)
			return "", text.RequestFailed[lang(ctx)]
		}
	} else if attempts > 1 {
		log.Printf("%q was handled after %d attempts", errMsg, attempts)
	}

	return response.Data[0].URL, response.Data[0].RevisedPrompt
}
