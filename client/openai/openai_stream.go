package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/types"
)

func streamOut(buffer *atomic.Value, isStreaming *atomic.Int64, channel chan<- string) {

	var previous string
	var completion string
	buffer.Store(previous)
	isStreaming.Store(1)

	for isStreaming.Load() == 1 {
		if completion != previous {
			channel <- completion + " ▌"
			previous = completion
		}
		time.Sleep(constants.RequestInterval)
		completion = buffer.Load().(string)
	}

	if len(completion) > 0 {
		channel <- completion
	}

	if isStreaming.Load() == 0 {
		close(channel)
	}
}

func streamIn(resp *http.Response, buffer *atomic.Value) (*types.OpenAIResponse, error) {

	var (
		err      error
		bts      []byte
		builder  strings.Builder
		response = &types.OpenAIResponse{}
		reader   = bufio.NewReader(resp.Body)
		prefix   = []byte{'d', 'a', 't', 'a', ':', ' '}
		userID   = resp.Request.Context().Value("user_id").(int64)
	)

	response.Choices = append(response.Choices, types.Choice{})

	for {
		bts, err = reader.ReadBytes('\n')
		if err != nil {
			log.Printf("%v", err)
			return nil, err
		}

		if bts[0] == '\n' {
			continue
		}

		bts = bytes.TrimPrefix(bts, prefix)
		if string(bts) == "[DONE]" {
			response.Choices[0].FinishReason = "[DONE]"
			break
		}

		err = json.Unmarshal(bts, &response)
		if err != nil {
			log.Printf("can't decode response for %d: %v\nbody: %s", userID, err, string(bts))
			return nil, fmt.Errorf("can't decode response for %d", userID)
		}

		if response.Choices[0].FinishReason != "" {
			if response.Choices[0].FinishReason != "stop" {
				log.Printf("finish reason for %d isn't \"stop\": %q", userID, response.Choices[0].FinishReason)
			}
			break
		}

		builder.WriteString(response.Choices[0].Delta.Content)
		buffer.Store(builder.String())
	}

	if len(builder.String()) == 0 {
		log.Printf("empty completion for %d", userID)
		return nil, fmt.Errorf("empty completion for %d", userID)
	}

	response.Choices[0].Message.Role = "assistant"
	response.Choices[0].Message.Content = builder.String()

	return response, nil
}
