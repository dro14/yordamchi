package service

import (
	"context"
	"errors"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ErrUnsupportedFormat = errors.New("invalid file")

func (s *Service) Load(ctx context.Context, document *tgbotapi.Document) error {
	request := map[string]any{
		"file_id":   document.FileID,
		"file_name": document.FileName,
		"user_id":   ctx.Value("user_id").(int64),
	}
	response, err := s.makeRequest(ctx, request, s.baseURL+"load")
	if err != nil {
		return err
	}
	if response["success"] == false {
		log.Println("can't load file:", response["error"])
		return ErrUnsupportedFormat
	}
	return nil
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
	return response["file_name"].(string)
}

func (s *Service) Clear(ctx context.Context) error {
	request := map[string]any{
		"user_id": ctx.Value("user_id").(int64),
	}
	_, err := s.makeRequest(ctx, request, s.baseURL+"clear")
	if err != nil {
		return err
	}
	return nil
}
