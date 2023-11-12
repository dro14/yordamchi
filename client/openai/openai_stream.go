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

	for isStreaming.Load() == 1 {
		completion = buffer.Load().(string)
		if completion != previous {
			channel <- completion + " ▌"
			previous = completion
		}
		time.Sleep(constants.RequestInterval)
	}

	completion = buffer.Load().(string)
	if len(strings.TrimSpace(completion)) > 0 {
		channel <- completion
	}
	if isStreaming.Load() == 0 {
		close(channel)
	}
}

func streamIn(resp *http.Response, buffer *atomic.Value) (*types.Response, error) {
	var builder strings.Builder
	response := &types.Response{}
	reader := bufio.NewReader(resp.Body)

	for {
		bts, err := reader.ReadBytes('\n')
		if err != nil {
			if strings.HasPrefix(err.Error(), "stream error") {
				return nil, fmt.Errorf("stream error for %s", id(resp.Request.Context()))
			}
			log.Printf("%s", err)
			return nil, err
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
			log.Printf("can't decode response for %s: %s\nbody: %s", id(resp.Request.Context()), err, string(bts))
			return nil, fmt.Errorf("can't decode response for %s", id(resp.Request.Context()))
		}

		if len(response.Choices) == 0 {
			log.Printf("empty choices for %s", id(resp.Request.Context()))
			continue
		}

		if response.Choices[0].FinishReason != "" {
			if response.Choices[0].FinishReason != "stop" {
				log.Printf("finish reason for %s isn't \"stop\": %q", id(resp.Request.Context()), response.Choices[0].FinishReason)
			}
			break
		} else if response.Choices[0].FinishDetails.Type != "" {
			if response.Choices[0].FinishDetails.Type != "stop" {
				log.Printf("finish details type for %s isn't \"stop\": %q", id(resp.Request.Context()), response.Choices[0].FinishDetails.Type)
			}
			break
		}

		builder.WriteString(response.Choices[0].Delta.Content)
		buffer.Store(builder.String())
	}

	if len(strings.TrimSpace(builder.String())) == 0 {
		return nil, fmt.Errorf("empty completion for %s", id(resp.Request.Context()))
	}
	response.Choices[0].Message.Role = "assistant"
	response.Choices[0].Message.Content = builder.String()
	return response, nil
}
