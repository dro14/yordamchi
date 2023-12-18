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
	ErrNotFound      = errors.New("400 Bad Request: message to edit/delete not found")
	ErrNotModified   = errors.New("400 Bad Request: message is not modified")
	ErrCantDelete    = errors.New("400 Bad Request: message can't be deleted for everyone")
	ErrMarkdown      = errors.New("400 Bad Request: can't parse entities")
	ErrForbidden     = errors.New("403 Forbidden: bot was blocked by the user")
	ErrRequestFailed = errors.New("request failed")
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
				return nil, ErrNotFound
			case is("message is not modified"):
				return nil, ErrNotModified
			case is("message can't be deleted for everyone"):
				return nil, ErrCantDelete
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
			default:
				log.Printf("user %d: failed after %d attempts: %d %s", id(ctx), attempts, resp.ErrorCode, resp.Description)
				return nil, ErrRequestFailed
			}
		} else {
			log.Printf("user %d: can't make request: %s", id(ctx), err)
			return nil, err
		}
	}
	if attempts > 1 {
		log.Printf("user %d: request was handled after %d attempts", id(ctx), attempts)
	}
	return resp, nil
}
