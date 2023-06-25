package telegram

import (
	"context"
	"fmt"
	"github.com/dro14/yordamchi/lib/constants"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/gotd/td/tg"
)

type Client struct {
	bot *tg.Client
}

func New(bot *tg.Client) *Client {
	return &Client{
		bot: bot,
	}
}

func (c *Client) SendMessage(ctx context.Context, message string, replyToMsgID int, keyboard *tg.ReplyInlineMarkup) (int, error) {

	userID := ctx.Value("user_id").(int64)

	request := &tg.MessagesSendMessageRequest{
		Peer:         &tg.InputPeerUser{UserID: userID},
		Message:      message,
		ReplyToMsgID: replyToMsgID,
		RandomID:     time.Now().UnixNano(),
		NoWebpage:    true,
	}

	if keyboard != nil {
		request.ReplyMarkup = keyboard
	}

	retryDelay := constants.RetryDelay
	attempts := 0
Retry:
	attempts++
	resp, err := c.bot.MessagesSendMessage(ctx, request)
	if err != nil {

		log.Printf("can't send message to %d: %v", userID, err)

		switch {
		case strings.Contains(err.Error(), e.BrokenPipe):
			log.Fatalf("fatal error: restarting bot")
		case strings.Contains(err.Error(), e.UserBlocked):
			return 0, e.UserBlockedError
		case strings.Contains(err.Error(), e.MessageEmpty),
			strings.Contains(err.Error(), e.MessageTooLong):
			log.Printf("%q", request.Message)
			return 0, err
		case strings.Contains(err.Error(), e.TooManyRequests):
			_, str, _ := strings.Cut(err.Error(), e.TooManyRequests)
			str = str[1 : len(str)-1]
			seconds, _ := strconv.Atoi(str)
			retryDelay = time.Duration(seconds) * time.Second
		}

		if attempts < constants.RetryAttempts {
			functions.Sleep(&retryDelay)
			goto Retry
		}
		return 0, err
	}

	response, ok := resp.(*tg.UpdateShortSentMessage)
	if !ok {
		log.Printf("can't decode response for %d: %v", userID, err)
		return 0, fmt.Errorf("can't decode response for %d", userID)
	}

	if attempts > 1 {
		log.Printf("sending message to %d was handled after %d attempts", userID, attempts)
	}

	return response.ID, nil
}

func (c *Client) EditMessage(ctx context.Context, message string, messageID int, keyboard *tg.ReplyInlineMarkup) error {

	userID := ctx.Value("user_id").(int64)

	request := &tg.MessagesEditMessageRequest{
		Peer:      &tg.InputPeerUser{UserID: userID},
		ID:        messageID,
		NoWebpage: true,
	}

	if len(message) > 0 {
		request.Message = message
	}

	if keyboard != nil {
		request.ReplyMarkup = keyboard
	}

	retryDelay := constants.RetryDelay
	attempts := 0
Retry:
	attempts++
	_, err := c.bot.MessagesEditMessage(ctx, request)
	if err != nil {

		log.Printf("can't edit message for %d: %v", userID, err)

		switch {
		case strings.Contains(err.Error(), e.BrokenPipe):
			log.Fatalf("fatal error: restarting bot")
		case strings.Contains(err.Error(), e.UserBlocked):
			return e.UserBlockedError
		case strings.Contains(err.Error(), e.MessageNotFound):
			return e.UserDeletedMessage
		case strings.Contains(err.Error(), e.MessageEmpty),
			strings.Contains(err.Error(), e.MessageTooLong),
			strings.Contains(err.Error(), e.MessageNotModified):
			log.Printf("%q", request.Message)
			return err
		case strings.Contains(err.Error(), e.TooManyRequests):
			_, str, _ := strings.Cut(err.Error(), e.TooManyRequests)
			str = str[1 : len(str)-1]
			seconds, _ := strconv.Atoi(str)
			retryDelay = time.Duration(seconds) * time.Second
		}

		if attempts < constants.RetryAttempts {
			functions.Sleep(&retryDelay)
			goto Retry
		}
		return err
	}

	if attempts > 1 {
		log.Printf("editing message for %d was handled after %d attempts", userID, attempts)
	}

	return nil
}

func (c *Client) SetTyping(ctx context.Context, isTyping *atomic.Bool) {

	userID := ctx.Value("user_id").(int64)

	request := &tg.MessagesSetTypingRequest{
		Peer:   &tg.InputPeerUser{UserID: userID},
		Action: &tg.SendMessageTypingAction{},
	}

Loop:
	for isTyping.Load() {
		_, err := c.bot.MessagesSetTyping(ctx, request)
		if err != nil {

			log.Printf("can't set typing for %d: %v", userID, err)

			switch {
			case strings.Contains(err.Error(), e.BrokenPipe):
				log.Fatalf("fatal error: restarting bot")
			case strings.Contains(err.Error(), e.UserBlocked):
				break Loop
			}
		}
		time.Sleep(5800 * time.Millisecond)
	}
}
