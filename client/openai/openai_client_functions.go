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
	resp, err := makeRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp, nil
	case http.StatusBadGateway, http.StatusServiceUnavailable:
		return nil, fmt.Errorf("%s for %s", resp.Status, id(ctx))
	default:
		bts, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		response := &types.Response{}
		err = json.Unmarshal(bts, response)
		if err != nil {
			log.Printf("%s for %s\ncan't decode response: %s\nbody: %q", resp.Status, id(ctx), err, string(bts))
		} else {
			log.Printf("%s for %s\ntype: %s\nmessage: %s", resp.Status, id(ctx), response.Error.Type, response.Error.Message)
		}
		return nil, fmt.Errorf("%s for %s", resp.Status, id(ctx))
	}
}

func makeRequest[T *types.Completions | *types.Generations](ctx context.Context, request T) (*http.Response, error) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Printf("can't encode request: %s", err)
		return nil, err
	}

	URL := ctx.Value("url").(string)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL, &buffer)
	if err != nil {
		log.Printf("can't create request: %s", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", keys[index])
	index++
	if index == len(keys) {
		index = 0
	}

	var client http.Client
	client.Timeout = 3 * time.Minute
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
