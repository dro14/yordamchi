package openai

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/dro14/yordamchi/storage/postgres"
	"github.com/dro14/yordamchi/storage/redis"
	"github.com/dro14/yordamchi/utils"
)

var template = map[string]string{
	"uz": "%s\n\nQUYIDA MAVZUGA OID MA'LUMOTLAR KELTIRILGAN. KERAK BO'LSA ULARDAN FOYDALAN.\n\n%s",
	"ru": "%s\n\nНИЖЕ ПРИВЕДЕНЫ СООТВЕТСТВУЮЩИЕ ТЕМЕ ФРАГМЕНТЫ ИНФОРМАЦИИ. ИСПОЛЬЗУЙ ИХ, ЕСЛИ ОНИ БУДУТ ПОЛЕЗНЫ.\n\n%s",
	"en": "%s\n\nTHE FOLLOWING ARE THE RELEVANT PIECES OF INFORMATION. USE THEM IF HELPFUL.\n\n%s",
}

func (o *OpenAI) ProcessCompletions(ctx context.Context, prompt string, msg *postgres.Message, channel chan<- string) {
	defer close(channel)
	defer utils.RecoverIfPanic()

	ctx, messages := o.redis.Context(ctx, &prompt)
	if userStatus(ctx) != redis.StatusFree && !strings.Contains(prompt, utils.Delim) {
		results := o.service.Search(ctx, prompt)
		if results != "" {
			if model(ctx) == models.GPT3 && lang(ctx) == "uz" {
				results = o.apis.Translate("auto", "en", results)
				messages[0].Content = fmt.Sprintf(template["en"], messages[0].Content, results)
			} else {
				messages[0].Content = fmt.Sprintf(template[lang(ctx)], messages[0].Content, results)
			}
		}
	}

	retryDelay := 10 * utils.RetryDelay
	var errMsg string
Retry:
	msg.Attempts++
	response, err := o.Completions(ctx, messages, channel)
	if err != nil {
		errMsg = err.Error()
		is := func(s string) bool {
			return strings.Contains(errMsg, s)
		}
		switch {
		case is("400 Bad Request"):
			channel <- text.BadRequest[lang(ctx)]
			return
		case is("stream error"):
			channel <- text.Error[lang(ctx)]
			retryDelay = 0
		case is("context deadline exceeded"),
			is("500 Internal Server Error"),
			is("502 Bad Gateway"):
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
	msg.CompletionLength = len([]rune(completion))

	o.redis.SetContext(ctx, prompt, completion)
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
		is := func(s string) bool {
			return strings.Contains(errMsg, s)
		}
		switch {
		case is("400 Bad Request"):
			return "", text.BadRequest[lang(ctx)]
		case is("context deadline exceeded"),
			is("500 Internal Server Error"),
			is("502 Bad Gateway"):
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

	URL := response.Data[0].URL
	path, _, found := strings.Cut(URL, ".png")
	if !found {
		log.Printf("user %s: can't find .png in %q", id(ctx), URL)
		return "", text.RequestFailed[lang(ctx)]
	}

	_, path, found = strings.Cut(path, "img-")
	if !found {
		log.Printf("user %s: can't find img- in %q", id(ctx), path)
		return "", text.RequestFailed[lang(ctx)]
	}
	path = "img-" + path + ".png"

	err = utils.DownloadFile(URL, path)
	if err != nil {
		log.Printf("user %s: can't download %q: %s", id(ctx), path, err)
		return "", text.RequestFailed[lang(ctx)]
	}

	return path, response.Data[0].RevisedPrompt
}
