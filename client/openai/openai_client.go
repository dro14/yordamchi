package openai

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/dro14/yordamchi/lib/types"
)

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

func Completion(ctx context.Context, messages []types.Message, maxTokens int, channel chan<- string) (*types.Response, error) {

	request := &types.Request{
		Model:     ctx.Value("model").(string) + "-0613",
		Messages:  messages,
		MaxTokens: maxTokens,
		Stream:    true,
		User:      fmt.Sprintf("%d", ctx.Value("user_id").(int64)),
	}

	resp, err := send(ctx, request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	buffer := &atomic.Value{}
	isStreaming := &atomic.Int64{}
	go streamOut(buffer, isStreaming, channel)

	response, err := streamIn(resp, buffer)
	if err != nil {
		isStreaming.Store(-1)
		return nil, err
	}

	isStreaming.Store(0)
	return response, nil
}
