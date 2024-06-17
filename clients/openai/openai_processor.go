package openai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/clients/openai/types"
	"github.com/dro14/yordamchi/processor/text"
	dbTypes "github.com/dro14/yordamchi/storage/postgres/types"
	"github.com/dro14/yordamchi/storage/redis/status"
	"github.com/dro14/yordamchi/utils"
)

func (o *OpenAI) ProcessCompletions(ctx context.Context, prompt string, msg *dbTypes.Message, channel chan<- string) {
	defer close(channel)
	defer utils.RecoverIfPanic()

	var completion, source string
	var tools []types.Tool
	messages := o.redis.Context(ctx, &prompt)
	if userStatus(ctx) != status.Free {
		source = o.service.Memory(ctx)
		if source == "GOOGLE" {
			tools = append(tools, googleSearch)
		} else {
			fileSearch.Function.Description = fileSearchDescription + source
			tools = append(tools, fileSearch)
		}
	}

	retryDelay := 10 * utils.RetryDelay
	var errMsg string
Retry:
	msg.Attempts++
	response, err := o.Completions(ctx, messages, tools, completion, channel)
	if err != nil {
		errMsg = err.Error()
		is := func(s string) bool {
			return strings.Contains(errMsg, s)
		}
		switch {
		case errors.Is(err, contextLengthExceeded):
			channel <- text.ContextLengthExceeded[lang(ctx)]
			return
		case errors.Is(err, inappropriateRequest):
			channel <- text.InappropriateRequest[lang(ctx)]
			return
		case errors.Is(err, badRequest):
			channel <- text.BadRequest[lang(ctx)]
			return
		case is("context deadline exceeded"), is("500 Internal Server Error"), is("502 Bad Gateway"):
			retryDelay = 0
		case is("stream error"):
			channel <- text.StreamError[lang(ctx)]
			retryDelay = 0
		}
		if msg.Attempts < utils.RetryAttempts {
			utils.Sleep(&retryDelay)
			goto Retry
		} else {
			log.Printf("%q failed after %d attempts", errMsg, msg.Attempts)
			channel <- text.FailedRequest[lang(ctx)]
			return
		}
	} else if msg.Attempts > 1 {
		if errMsg != "" {
			log.Printf("%q was handled after %d attempts", errMsg, msg.Attempts)
		}
	}

	msg.PromptTokens += response.Usage.PromptTokens
	msg.CompletionTokens += response.Usage.CompletionTokens

	if len(getToolCalls(response)) > 0 {
		var callResults []types.Message
		for _, toolCall := range getToolCalls(response) {
			var args map[string]string
			_ = json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
			query, ok := args["query"]
			if !ok {
				log.Printf("user %s: invalid JSON body from OpenAI: %s", id(ctx), toolCall.Function.Arguments)
				if msg.Attempts < utils.RetryAttempts {
					goto Retry
				} else {
					log.Printf("%q failed after %d attempts", errMsg, msg.Attempts)
					channel <- text.FailedRequest[lang(ctx)]
					return
				}
			}
			var result string
			if source == "GOOGLE" {
				result = o.service.GoogleSearch(ctx, query)
			} else {
				result = o.service.FileSearch(ctx, query)
			}
			callResults = append(callResults,
				types.Message{
					Role:       "tool",
					Content:    result,
					ToolCallID: toolCall.ID,
				},
			)
		}

		messages = append(messages, response.Choices[0].Message)
		messages = append(messages, callResults...)
		completion += getCompletion(response)
		if translate(ctx) {
			completion += fmt.Sprintf(text.Search["en"], source)
		} else {
			completion += fmt.Sprintf(text.Search[lang(ctx)], source)
		}

		if msg.Attempts < utils.RetryAttempts {
			goto Retry
		} else {
			log.Printf("%q failed after %d attempts", errMsg, msg.Attempts)
			channel <- text.FailedRequest[lang(ctx)]
			return
		}
	}

	if ctx.Value("stream") == true {
		completion = getCompletion(response)
	} else {
		completion += getCompletion(response)
	}
	msg.Output = completion
	msg.PromptLength += length(messages)
	msg.CompletionLength += len([]rune(completion))
	msg.FinishReason = getFinishReason(response)
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
		case errors.Is(err, contextLengthExceeded):
			return "", text.ContextLengthExceeded[lang(ctx)]
		case errors.Is(err, inappropriateRequest):
			return "", text.InappropriateRequest[lang(ctx)]
		case errors.Is(err, badRequest):
			return "", text.BadRequest[lang(ctx)]
		case is("context deadline exceeded"), is("500 Internal Server Error"), is("502 Bad Gateway"):
			retryDelay = 0
		}
		if attempts < utils.RetryAttempts {
			utils.Sleep(&retryDelay)
			goto Retry
		} else {
			log.Printf("%q failed after %d attempts", errMsg, attempts)
			return "", text.FailedRequest[lang(ctx)]
		}
	} else if attempts > 1 {
		log.Printf("%q was handled after %d attempts", errMsg, attempts)
	}

	URL := response.Data[0].URL
	path, _, found := strings.Cut(URL, ".png")
	if !found {
		log.Printf("user %s: can't find .png in %q", id(ctx), URL)
		return "", text.FailedRequest[lang(ctx)]
	}

	_, path, found = strings.Cut(path, "img-")
	if !found {
		log.Printf("user %s: can't find img- in %q", id(ctx), path)
		return "", text.FailedRequest[lang(ctx)]
	}
	path = "img-" + path + ".png"

	err = utils.DownloadFile(URL, path)
	if err != nil {
		log.Printf("user %s: can't download %q: %s", id(ctx), path, err)
		return "", text.FailedRequest[lang(ctx)]
	}

	return path, response.Data[0].RevisedPrompt
}

var template = `IN THE LANGUAGE OF CONVERSATION, GENERATE 3 VERY BRIEF FOLLOW-UP QUESTIONS THAT THE USER WOULD LIKELY ASK NEXT.
RESPOND IN JSON FORMAT:
{"questions":[<string>,<string>,<string>]}`

func (o *OpenAI) ProcessFollowUps(ctx context.Context) []string {
	ctx = context.WithValue(ctx, "stream", false)
	ctx = context.WithValue(ctx, "json_mode", true)
	ctx = context.WithValue(ctx, "translate", false)
	ctx = context.WithValue(ctx, "user_status", o.redis.UserStatus(ctx))
	ctx = context.WithValue(ctx, "model", models.GPT3)

	if userStatus(ctx) == status.Exhausted {
		return text.DefaultQuestions[lang(ctx)]
	}

	messages := []types.Message{{Role: "system", Content: template}}
	messages = append(messages, o.redis.Messages(ctx)...)

	response, err := o.Completions(ctx, messages, nil, "", make(chan string, 1))
	if err != nil {
		log.Printf("user %s: can't generate follow-up questions: %v", id(ctx), err)
		return text.DefaultQuestions[lang(ctx)]
	}

	bts := []byte(response.Choices[0].Message.Content.(string))
	var questions map[string][]string
	err = json.Unmarshal(bts, &questions)
	if err != nil {
		log.Printf("user %s: can't decode response: %v\nbody: %s", id(ctx), err, bts)
		return text.DefaultQuestions[lang(ctx)]
	}

	if userStatus(ctx) != status.Premium && lang(ctx) == "uz" {
		for i, question := range questions["questions"] {
			questions["questions"][i] = o.apis.Translate("en", "uz", question)
		}
	}

	return questions["questions"]
}
