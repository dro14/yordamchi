package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dro14/yordamchi/lib/types"
)

func send[T *types.Completions | *types.Generations](ctx context.Context, request T) (*http.Response, error) {
	userID := ctx.Value("user_id").(int64)

	resp, err := makeRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp, nil
	case http.StatusBadGateway, http.StatusServiceUnavailable:
		return nil, fmt.Errorf("%s for %d", resp.Status, userID)
	default:
		bts, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		var response types.Response
		err = json.Unmarshal(bts, &response)
		if err != nil {
			log.Printf("%s for %d\ncan't decode response: %v\nbody: %q", resp.Status, userID, err, string(bts))
		} else {
			log.Printf("%s for %d\ntype: %s\nmessage: %s", resp.Status, userID, response.Error.Type, response.Error.Message)
		}
		return nil, fmt.Errorf("%s for %d", resp.Status, userID)
	}
}

func makeRequest[T *types.Completions | *types.Generations](ctx context.Context, request T) (*http.Response, error) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Printf("can't encode request: %v", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ctx.Value("url").(string), &buffer)
	if err != nil {
		log.Printf("can't create request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", keys[index])

	index++
	if index == len(keys) {
		index = 0
	}

	var client http.Client
	client.Timeout = 10 * time.Minute

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
