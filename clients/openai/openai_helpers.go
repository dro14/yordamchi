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

const (
	googleSearchDescription = "Searches for real-time information in Google"
	fileSearchDescription   = "Searches for additional information in the file the user provided. File name: "
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

	choice := types.Choice{}
	response := &types.Response{}
	reader := bufio.NewReader(resp.Body)
	var previous, toolCallID string
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
			break
		}

		err = json.Unmarshal(bts, response)
		if err != nil {
			log.Printf("user %s: can't decode response: %s\nbody: %s", id(ctx), err, bts)
			return nil, fmt.Errorf("user %s: can't decode response", id(ctx))
		}

		if len(response.Choices) == 0 {
			continue
		}

		if response.Choices[0].Delta.Role != "" {
			choice.Message.Role = response.Choices[0].Delta.Role
		}

		if response.Choices[0].FinishReason != "" {
			choice.FinishReason = response.Choices[0].FinishReason
		}

		if response.Choices[0].Delta.ToolCalls != nil {
			if getToolCallID(response) != "" && getToolCallID(response) != toolCallID {
				toolCallID = getToolCallID(response)
				if choice.Message.ToolCalls != nil {
					choice.Message.ToolCalls[len(choice.Message.ToolCalls)-1].Function.Arguments = args.String()
					args.Reset()
				}
				choice.Message.ToolCalls = append(choice.Message.ToolCalls, response.Choices[0].Delta.ToolCalls[0])
			}
			args.WriteString(response.Choices[0].Delta.ToolCalls[0].Function.Arguments)
			response.Choices[0].Delta.ToolCalls = nil
		}

		content.WriteString(response.Choices[0].Delta.Content)
		response.Choices[0].Delta.Content = ""
		completion = strings.TrimSpace(content.String())
		if send.Load() && completion != previous {
			channel <- completion + " ▌"
			previous = completion
			send.Store(false)
		}
	}

	stream.Store(false)
	choice.Message.Content = completion
	if choice.Message.ToolCalls != nil {
		choice.Message.ToolCalls[len(choice.Message.ToolCalls)-1].Function.Arguments = args.String()
	}
	response.Choices = append(response.Choices, choice)
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
	for _, message := range messages {
		promptLength += len([]rune(fmt.Sprintf("<|start|>%v\n%v<|end|>\n", message.Role, message.Content)))
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

func translate(ctx context.Context) bool {
	return ctx.Value("translate").(bool)
}

func getCompletion(response *types.Response) string {
	content, ok := response.Choices[0].Message.Content.(string)
	if !ok {
		return ""
	}
	return content
}

func getFinishReason(response *types.Response) string {
	return response.Choices[0].FinishReason
}

func getToolCalls(response *types.Response) []types.ToolCall {
	return response.Choices[0].Message.ToolCalls
}

func getToolCallID(response *types.Response) string {
	return response.Choices[0].Delta.ToolCalls[0].ID
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
