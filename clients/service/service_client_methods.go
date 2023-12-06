package service

import (
	"context"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) Load(ctx context.Context, document *tgbotapi.Document) string {
	request := map[string]any{
		"file_id":   document.FileID,
		"file_name": document.FileName,
		"user_id":   ctx.Value("user_id").(int64),
	}
	response, err := s.makeRequest(ctx, request, s.baseURL+"load")
	if err != nil {
		return err.Error()
	}
	if response["success"] == false {
		errMsg, supported, _ := strings.Cut(response["error"].(string), "\n")
		log.Println("can't load file:", errMsg)
		return supported
	}
	return ""
}

func (s *Service) Search(ctx context.Context, query string) string {
	request := map[string]any{
		"query":   query,
		"lang":    ctx.Value("language_code").(string),
		"user_id": ctx.Value("user_id").(int64),
	}
	response, err := s.makeRequest(ctx, request, s.baseURL+"search")
	if err != nil {
		return ""
	}
	return response["results"].(string)
}

func (s *Service) Memory(ctx context.Context) string {
	request := map[string]any{
		"user_id": ctx.Value("user_id").(int64),
	}
	response, err := s.makeRequest(ctx, request, s.baseURL+"memory")
	if err != nil {
		return ""
	}
	return response["source"].(string)
}

func (s *Service) Delete(ctx context.Context) error {
	request := map[string]any{
		"user_id": ctx.Value("user_id").(int64),
	}
	_, err := s.makeRequest(ctx, request, s.baseURL+"delete")
	if err != nil {
		return err
	}
	return nil
}
