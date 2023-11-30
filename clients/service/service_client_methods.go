package service

import (
	"context"
)

func (s *Service) Search(ctx context.Context, query string) string {
	request := map[string]any{
		"query":   query,
		"lang":    ctx.Value("language_code").(string),
		"user_id": ctx.Value("user_id").(int64),
	}

	response, err := s.makeRequest(ctx, request, s.baseURL+"/search")
	if err != nil {
		return ""
	}

	if response["success"] == false {
		return ""
	}

	return response["results"].(string)
}
