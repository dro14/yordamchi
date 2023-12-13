package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (s *Service) makeRequest(ctx context.Context, request map[string]any, url string) (map[string]any, error) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Println("can't encode request:", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buffer)
	if err != nil {
		log.Println("can't create request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("can't do request:", err)
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("can't read response:", err)
		return nil, err
	}

	response := make(map[string]any)
	err = json.Unmarshal(bts, &response)
	if err != nil {
		log.Printf("can't decode response: %s\nbody: %s", err, bts)
		return nil, err
	}

	return response, nil
}
