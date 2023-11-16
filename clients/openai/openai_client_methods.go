package openai

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/dro14/yordamchi/clients/openai/types"
)

const Baseurl = "https://api.openai.com/v1/"
const ChatCompletions = "chat/completions"
const ImagesGenerations = "images/generations"

func (o *OpenAI) Completions(ctx context.Context, messages []types.Message, channel chan<- string) (*types.Response, error) {
	ctx = context.WithValue(ctx, "url", Baseurl+ChatCompletions)
	request := &types.Completions{
		Model:     ctx.Value("model").(string),
		Messages:  messages,
		MaxTokens: 4096,
		Stream:    ctx.Value("stream").(bool),
		User:      id(ctx),
	}

	resp, err := o.send(ctx, request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	response := &types.Response{}
	if ctx.Value("stream") == true {
		response, err = streamResponse(ctx, resp, channel)
	} else {
		response, err = decodeResponse(ctx, resp)
	}
	if err != nil {
		return nil, err
	} else if len(strings.TrimSpace(response.Choices[0].Message.Content.(string))) == 0 {
		return nil, fmt.Errorf("user %s: empty completion", id(ctx))
	} else if response.Choices[0].FinishReason != "stop" {
		log.Printf("user %s: finish reason isn't \"stop\": %q", id(ctx), response.Choices[0].FinishReason)
	}
	return response, nil
}

func (o *OpenAI) Generations(ctx context.Context, prompt string) (*types.Response, error) {
	ctx = context.WithValue(ctx, "url", Baseurl+ImagesGenerations)
	request := &types.Generations{
		Prompt:  prompt,
		Model:   "dall-e-3",
		Quality: "hd",
		Size:    "1024x1024",
		Style:   "vivid",
		User:    id(ctx),
	}

	resp, err := o.send(ctx, request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	response, err := decodeResponse(ctx, resp)
	if err != nil {
		return nil, err
	}
	return response, nil
}
