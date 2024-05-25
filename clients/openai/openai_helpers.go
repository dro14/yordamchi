package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dro14/yordamchi/clients/openai/types"
	"github.com/dro14/yordamchi/storage/redis"
	"github.com/dro14/yordamchi/utils"
)

func streamResponse(ctx context.Context, resp *http.Response, completion string, channel chan<- string) (*types.Response, error) {
	var stream, send atomic.Bool
	go func() {
		stream.Store(true)
		for stream.Load() {
			send.Store(true)
			time.Sleep(utils.ReqInterval)
		}
	}()

	response := &types.Response{Choices: []types.Choice{{}}}
	reader := bufio.NewReader(resp.Body)
	var previous string
	var content, args strings.Builder
	content.WriteString(completion)

	for {
		bts, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("user %s: %s", id(ctx), err)
			return nil, fmt.Errorf("user %s: stream error", id(ctx))
		}

		if bts[0] == '\n' {
			continue
		}

		bts = bytes.TrimPrefix(bts, []byte("data: "))
		if string(bts) == "[DONE]\n" {
			response.Choices[0].FinishReason = "done"
			break
		}

		err = json.Unmarshal(bts, response)
		if err != nil {
			log.Printf("user %s: can't decode response: %s\nbody: %s", id(ctx), err, bts)
			return nil, fmt.Errorf("user %s: can't decode response", id(ctx))
		}

		if getFinishReason(response) != "" {
			break
		} else if response.Choices[0].FinishDetails.Type != "" {
			response.Choices[0].FinishReason = response.Choices[0].FinishDetails.Type
			break
		}

		if response.Choices[0].Delta.ToolCalls != nil {
			if response.Choices[0].Message.ToolCalls == nil {
				response.Choices[0].Message.ToolCalls = response.Choices[0].Delta.ToolCalls
			}
			args.WriteString(response.Choices[0].Delta.ToolCalls[0].Function.Arguments)
			response.Choices[0].Delta.Content = ""
		}

		content.WriteString(response.Choices[0].Delta.Content)
		response.Choices[0].Delta.ToolCalls = nil
		completion = strings.TrimSpace(content.String())
		if send.Load() && completion != previous {
			channel <- completion + " ▌"
			previous = completion
			send.Store(false)
		}
	}

	stream.Store(false)
	response.Choices[0].Message.Role = response.Choices[0].Delta.Role
	response.Choices[0].Message.Content = completion
	if response.Choices[0].Message.ToolCalls != nil {
		response.Choices[0].Message.ToolCalls[0].Function.Arguments = args.String()
	}
	return response, nil
}

func decodeResponse(ctx context.Context, resp *http.Response) (*types.Response, error) {
	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("user %s: can't read response: %s", id(ctx), err)
		return nil, fmt.Errorf("user %s: can't read response", id(ctx))
	}

	response := &types.Response{}
	err = json.Unmarshal(bts, response)
	if err != nil {
		log.Printf("user %s: can't decode response: %s\nbody: %s", id(ctx), err, bts)
		return nil, fmt.Errorf("user %s: can't decode response", id(ctx))
	}
	return response, nil
}

func length(messages []types.Message) int {
	var promptLength int
	for i := range messages {
		promptLength += len([]rune(fmt.Sprintf("<|start|>%s\n%v<|end|>\n", messages[i].Role, messages[i].Content)))
	}
	return promptLength
}

func id(ctx context.Context) string {
	return fmt.Sprintf("%d", ctx.Value("user_id").(int64))
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

func userStatus(ctx context.Context) redis.UserStatus {
	return ctx.Value("user_status").(redis.UserStatus)
}

func getContent(response *types.Response) string {
	if response.Choices[0].Message.Content == nil {
		return ""
	}
	return response.Choices[0].Message.Content.(string)
}

func getArgs(response *types.Response) string {
	if response.Choices[0].Message.ToolCalls == nil {
		return ""
	}
	return response.Choices[0].Message.ToolCalls[0].Function.Arguments
}

func getFinishReason(response *types.Response) string {
	return response.Choices[0].FinishReason
}

var googleSearch = types.Tool{
	Type: "function",
	Function: types.Function{
		Name: "google_search",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{"type": "string"},
			},
		},
	},
}

var fileSearch = types.Tool{
	Type: "function",
	Function: types.Function{
		Name: "file_search",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{"type": "string"},
			},
		},
	},
}
