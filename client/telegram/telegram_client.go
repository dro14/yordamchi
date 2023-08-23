package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	cache "github.com/gotd/contrib/redis"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

var (
	bot *tgbotapi.BotAPI
	api *tg.Client
)

func Init() {

	token, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatalf("bot token is not specified")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("can't initialize bot: %v", err)
	}

	sessionStorage := cache.NewSessionStorage(redis.Client, "main_bot_session")
	done := make(chan bool)

	go func() {
		if err = telegram.BotFromEnvironment(
			context.Background(),
			telegram.Options{SessionStorage: sessionStorage},
			func(ctx context.Context, client *telegram.Client) error {
				api = client.API()
				done <- true
				return nil
			},
			telegram.RunUntilCanceled,
		); err != nil {
			log.Fatalf("can't connect client: %v", err)
		}
	}()

	<-done
}

func SendMessage(ctx context.Context, message string, replyToMsgID int, keyboard *tg.ReplyInlineMarkup) (int, error) {

	userID := ctx.Value("user_id").(int64)

	request := &tg.MessagesSendMessageRequest{
		Peer:      &tg.InputPeerUser{UserID: userID},
		Message:   message,
		ReplyTo:   &tg.InputReplyToMessage{ReplyToMsgID: replyToMsgID},
		RandomID:  time.Now().UnixNano(),
		NoWebpage: true,
	}

	if keyboard != nil {
		request.ReplyMarkup = keyboard
	}

	retryDelay := constants.RetryDelay
	attempts := 0
Retry:
	attempts++
	resp, err := api.MessagesSendMessage(ctx, request)
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

func EditMessage(ctx context.Context, message string, messageID int, keyboard *tg.ReplyInlineMarkup) error {

	userID := ctx.Value("user_id").(int64)

	request := &tg.MessagesEditMessageRequest{
		Peer:      &tg.InputPeerUser{UserID: userID},
		Message:   message,
		ID:        messageID,
		NoWebpage: true,
	}

	if keyboard != nil {
		request.ReplyMarkup = keyboard
	}

	retryDelay := constants.RetryDelay
	attempts := 0
Retry:
	attempts++
	_, err := api.MessagesEditMessage(ctx, request)
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

func SetTyping(ctx context.Context, isTyping *atomic.Bool) {

	userID := ctx.Value("user_id").(int64)

	request := &tg.MessagesSetTypingRequest{
		Peer:   &tg.InputPeerUser{UserID: userID},
		Action: &tg.SendMessageTypingAction{},
	}

Loop:
	for isTyping.Load() {
		_, err := api.MessagesSetTyping(ctx, request)
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

func Send(ctx context.Context, message string, addButton bool) error {

	userID := ctx.Value("user_id").(int64)
	lang := ctx.Value("language_code").(string)
	message = functions.MarkdownV2(message)

	config := tgbotapi.NewMessage(userID, message)
	config.ParseMode = tgbotapi.ModeMarkdownV2
	config.DisableWebPagePreview = true
	if addButton {
		config.ReplyMarkup = newChatButton(lang)
	}

	_, err := bot.Request(config)
	if err != nil {
		log.Printf("can't send message to %d: %v", userID, err)
		return err
	}

	return nil
}

func Edit(ctx context.Context, message string, messageID int, addButton bool) error {

	userID := ctx.Value("user_id").(int64)
	lang := ctx.Value("language_code").(string)
	message = functions.MarkdownV2(message)

	config := tgbotapi.NewEditMessageText(userID, messageID, message)
	config.ParseMode = tgbotapi.ModeMarkdownV2
	config.DisableWebPagePreview = true
	if addButton {
		config.ReplyMarkup = newChatButton(lang)
	}

	_, err := bot.Request(config)
	if err != nil {
		log.Printf("can't edit message for %d: %v", userID, err)
		return err
	}

	return nil
}

func newChatButton(lang string) *tgbotapi.InlineKeyboardMarkup {

	text := map[string]string{
		"uz": "💬 Yangi suhbat 💬",
		"ru": "💬 Новый разговор 💬",
		"en": "💬 New chat 💬",
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				text[lang], "new_chat",
			),
		),
	)

	return &keyboard
}
