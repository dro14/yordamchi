package service

import (
	"context"
	"log"
	"strings"

	"github.com/dro14/yordamchi/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s *Service) Load(ctx context.Context, document *tgbotapi.Document) string {
	request := map[string]any{
		"file_id":   document.FileID,
		"file_name": document.FileName,
		"user_id":   id(ctx),
	}
	response, err := s.makeRequest(ctx, request, s.baseURL+"load")
	if err != nil {
		return err.Error()
	}
	if response["success"] == false {
		errMsg, supported, found := strings.Cut(response["error"].(string), utils.Delim)
		log.Printf("user %d: can't load file: %s", id(ctx), errMsg)
		if found {
			return supported
		} else {
			return errMsg
		}
	}
	return ""
}

func (s *Service) Search(ctx context.Context, query string) string {
	request := map[string]any{
		"query":   query,
		"user_id": id(ctx),
	}
	response, err := s.makeRequest(ctx, request, s.baseURL+"search")
	if err != nil && response["success"] == false {
		return ""
	}
	return response["results"].(string)
}

func (s *Service) GoogleSearch(ctx context.Context, query string) string {
	request := map[string]any{
		"query": query,
		"lang":  lang(ctx),
	}
	response, err := s.makeRequest(ctx, request, s.baseURL+"google_search")
	if err != nil {
		return ""
	}
	if response["success"] == false {
		log.Printf("user %d: can't search google: %s", id(ctx), response["error"])
		return ""
	}
	return response["results"].(string)
}

func (s *Service) Memory(ctx context.Context) string {
	request := map[string]any{"user_id": id(ctx)}
	response, err := s.makeRequest(ctx, request, s.baseURL+"memory")
	if err != nil {
		return ""
	}
	return response["source"].(string)
}

func (s *Service) Delete(ctx context.Context) {
	request := map[string]any{"user_id": id(ctx)}
	response, err := s.makeRequest(ctx, request, s.baseURL+"delete")
	if err != nil {
		return
	}
	if response["success"] == false {
		log.Printf("user %d: can't delete file: %s", id(ctx), response["error"])
	}
}

func (s *Service) Logs(ctx context.Context) {
	request := map[string]any{"user_id": id(ctx)}
	response, err := s.makeRequest(ctx, request, s.baseURL+"logs")
	if err != nil {
		return
	}
	if response["success"] == false {
		log.Printf("user %d: can't get logs: %s", id(ctx), response["error"])
	}
}
