package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dro14/yordamchi/clients/openai/types"
	"io"
	"log"
	"net/http"
	"time"
)

func (o *OpenAI) send(ctx context.Context, request any) (*http.Response, error) {
	resp, err := o.makeRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp, nil
	case http.StatusBadGateway, http.StatusServiceUnavailable:
		return nil, fmt.Errorf("user %s: %s", id(ctx), resp.Status)
	default:
		bts, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		response := &types.Response{}
		err = json.Unmarshal(bts, response)
		if err != nil {
			log.Printf("user %s: %s\ncan't decode response: %s\nbody: %s", id(ctx), resp.Status, err, bts)
		} else {
			log.Printf("user %s: %s\ntype: %s\nmessage: %s", id(ctx), resp.Status, response.Error.Type, response.Error.Message)
		}
		return nil, fmt.Errorf("user %s: %s", id(ctx), resp.Status)
	}
}

func (o *OpenAI) makeRequest(ctx context.Context, request any) (*http.Response, error) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Println("can't encode request:", err)
		return nil, err
	}

	URL := ctx.Value("url").(string)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, URL, &buffer)
	if err != nil {
		log.Println("can't create request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", o.keys[o.index])
	o.index++
	if o.index == len(o.keys) {
		o.index = 0
	}

	var client http.Client
	client.Timeout = 5 * time.Minute
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
