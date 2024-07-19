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

func (s *Service) FileSearch(ctx context.Context, query string) string {
	request := map[string]any{
		"query":   query,
		"user_id": id(ctx),
	}
	response, err := s.makeRequest(ctx, request, s.baseURL+"file_search")
	if err != nil || response["success"] == false {
		return "no results"
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
		return "no results"
	}
	if response["success"] == false {
		log.Printf("user %d: can't search google: %s", id(ctx), response["error"])
		return "no results"
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

func (s *Service) Files(ctx context.Context) string {
	response, err := s.makeRequest(ctx, nil, s.baseURL+"files")
	if err != nil {
		return "🚨 ERROR 🚨"
	}
	if response["success"] == false {
		log.Printf("user %d: can't get files: %s", id(ctx), response["error"])
		return response["error"].(string)
	}
	return "*Uploaded files:*\n\n" + response["files"].(string)
}

func (s *Service) Latex2Text(ctx context.Context, str string) string {
	LaTeXes := utils.LaTeXRgx.FindAllStringSubmatch(str, -1)
	if len(LaTeXes) == 0 {
		return str
	}
	var matches, latex []string
	for _, ltx := range LaTeXes {
		matches = append(matches, ltx[0])
		latex = append(latex, preProcess(ltx[1]))
	}

	request := map[string]any{"latex": latex}
	response, err := s.makeRequest(ctx, request, s.baseURL+"latex2text")
	if err != nil {
		return str
	}

	text, ok := response["text"].([]any)
	if !ok {
		return str
	}

	for i, match := range matches {
		unicode := postProcess(text[i].(string))
		str = strings.Replace(str, match, "`"+unicode+"`", 1)
	}
	return str
}
