package service

import (
	"context"
	"log"
	"strings"

	"github.com/dro14/yordamchi/utils"
)

func (s *Service) GoogleSearch(ctx context.Context, query string) string {
	utils.SendInfoMessage(query)
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
	results := response["results"].(string)
	utils.SendInfoMessage(results)
	return results
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
