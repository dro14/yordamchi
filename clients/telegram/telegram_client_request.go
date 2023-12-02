package telegram

import (
	"context"
	"errors"
	"fmt"
	"github.com/dro14/yordamchi/utils"
	"log"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrMessageNotFound = errors.New("400 Bad Request: message to edit/delete not found")
	ErrMarkdown        = errors.New("400 Bad Request: can't parse entities")
	ErrForbidden       = errors.New("403 Forbidden: bot was blocked by the user")
	ErrRequestFailed   = errors.New("request failed")
)

func (t *Telegram) makeRequest(ctx context.Context, request tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	retryDelay := utils.RetryDelay
	attempts := 0
Retry:
	attempts++
	resp, err := t.bot.Request(request)
	if err != nil {
		if resp != nil {
			log.Printf("user %d: can't make request: %d %s", id(ctx), resp.ErrorCode, resp.Description)
			is := func(s string) bool {
				return strings.Contains(resp.Description, s)
			}
			switch {
			case is("Forbidden"):
				return nil, ErrForbidden
			case is("message") && is("not found"):
				return nil, ErrMessageNotFound
			case is("can't parse entities"):
				return nil, ErrMarkdown
			case is("Bad Request"):
				log.Printf("%+v", request)
				return nil, fmt.Errorf("400 %w", err)
			case is("Too Many Requests"):
				retryDelay = time.Duration(resp.Parameters.RetryAfter) * time.Second
				fallthrough
			case attempts < utils.RetryAttempts:
				utils.Sleep(&retryDelay)
				goto Retry
			}
		}
		log.Printf("user %d: failed after %d attempts: %d %s", id(ctx), attempts, resp.ErrorCode, resp.Description)
		return nil, ErrRequestFailed
	}
	return resp, nil
}
