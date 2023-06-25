package openai

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"

	"github.com/dro14/yordamchi/lib/types"
)

type Client struct {
	tokens []string
	index  int
}

func New() *Client {

	client := &Client{}

	for i := 0; ; i++ {
		key := fmt.Sprintf("OPENAI_TOKEN_%d", i)
		token, ok := os.LookupEnv(key)
		if !ok {
			break
		}
		client.tokens = append(client.tokens, "Bearer "+token)
	}

	if len(client.tokens) == 0 {
		log.Fatalf("openai token is not specified")
	}

	return client
}

func (c *Client) Completion(ctx context.Context, messages []types.Message, maxTokens int, channel chan<- string) (*types.Response, error) {

	request := &types.Request{
		Model:     ctx.Value("model").(string) + "-0613",
		Messages:  messages,
		MaxTokens: maxTokens,
		Stream:    true,
		User:      fmt.Sprintf("%d", ctx.Value("user_id").(int64)),
	}

	resp, err := c.send(ctx, request)
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
