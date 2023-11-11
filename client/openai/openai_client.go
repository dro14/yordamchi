package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"

	"github.com/dro14/yordamchi/lib/models"
	"github.com/dro14/yordamchi/lib/types"
)

const Baseurl = "https://api.openai.com/v1/"
const ChatCompletions = "chat/completions"
const ImagesGenerations = "images/generations"

var keys []string
var index int

func Init() {

	for i := 0; ; i++ {
		key := fmt.Sprintf("OPENAI_TOKEN_%d", i)
		token, ok := os.LookupEnv(key)
		if !ok {
			break
		}
		keys = append(keys, "Bearer "+token)
	}

	if len(keys) == 0 {
		log.Fatalf("openai token is not specified")
	}
}

func CompletionsWithStream(ctx context.Context, messages []types.Message, channel chan<- string) (*types.Response, error) {

	ctx = context.WithValue(ctx, "url", Baseurl+ChatCompletions)
	userID := ctx.Value("user_id").(int64)

	request := &types.Completions{
		Model:     ctx.Value("model").(string),
		Messages:  messages,
		MaxTokens: maxTokens(ctx),
		Stream:    true,
		User:      fmt.Sprintf("%d", userID),
	}

	if ctx.Value("model") == models.GPT4V {
		n := len(messages)
		URL, text, _ := strings.Cut(messages[n-1].Content, "\n\n\n")
		messages[n-1].Content = ""
		messages[n-1].Contents = []types.Content{
			{Type: "text", Text: text},
			{Type: "image_url", ImageURL: types.ImageURL{URL: URL}},
		}
		request.Messages = messages
	}

	resp, err := send(ctx, request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	buffer := &atomic.Value{}
	buffer.Store("")
	isStreaming := &atomic.Int64{}
	isStreaming.Store(1)
	go streamOut(buffer, isStreaming, channel)

	response, err := streamIn(resp, buffer)
	if err != nil {
		isStreaming.Store(-1)
		return nil, err
	}

	isStreaming.Store(0)
	return response, nil
}

func Completions(ctx context.Context, messages []types.Message) (*types.Response, error) {

	ctx = context.WithValue(ctx, "url", Baseurl+ChatCompletions)
	userID := ctx.Value("user_id").(int64)

	request := &types.Completions{
		Model:     ctx.Value("model").(string),
		Messages:  messages,
		MaxTokens: maxTokens(ctx),
		Stream:    false,
		User:      fmt.Sprintf("%d", userID),
	}

	resp, err := send(ctx, request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("can't read response for %d: %v", userID, err)
		return nil, fmt.Errorf("can't read response for %d", userID)
	}

	response := &types.Response{}
	err = json.Unmarshal(bts, response)
	if err != nil {
		log.Printf("can't decode response for %d: %v\nbody: %s", userID, err, string(bts))
		return nil, fmt.Errorf("can't decode response for %d", userID)
	}

	if response.Choices[0].FinishReason != "stop" {
		log.Printf("finish reason for %d isn't \"stop\": %q", userID, response.Choices[0].FinishReason)
	}

	if len(strings.TrimSpace(response.Choices[0].Message.Content)) == 0 {
		return nil, fmt.Errorf("empty completion for %d", userID)
	}

	return response, nil
}

func Generations(ctx context.Context, prompt string) string {

	ctx = context.WithValue(ctx, "url", Baseurl+ImagesGenerations)
	userID := ctx.Value("user_id").(int64)

	request := &types.Generations{
		Prompt:  prompt,
		Model:   "dall-e-3",
		Quality: "hd",
		Size:    "1024x1024",
		Style:   "vivid",
		User:    fmt.Sprintf("%d", userID),
	}

	resp, err := send(ctx, request)
	if err != nil {
		log.Printf("can't send request: %v", err)
		return ""
	}
	defer func() { _ = resp.Body.Close() }()

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("can't read response: %v", err)
		return ""
	}

	response := &types.Response{}
	err = json.Unmarshal(bts, response)
	if err != nil {
		log.Printf("can't decode response: %v\nbody: %s", err, string(bts))
		return ""
	}

	return response.Data[0].URL
}
