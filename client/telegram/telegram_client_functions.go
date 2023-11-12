package telegram

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Send(ctx context.Context, text string, replyToMsgID int, replyMarkup *tgbotapi.InlineKeyboardMarkup) (int, error) {
	formatted := functions.MarkdownV2(text)
	config := tgbotapi.NewMessage(id(ctx), formatted)
	config.ReplyToMessageID = replyToMsgID
	config.ReplyMarkup = replyMarkup
	config.DisableWebPagePreview = true
	config.ParseMode = tgbotapi.ModeMarkdownV2
	retryDelay := constants.RetryDelay
	attempts := 0
Retry:
	attempts++
	resp, err := bot.Request(config)
	if err != nil {
		description := resp.Description
		log.Printf("can't send message to %d: %d %s", id(ctx), resp.ErrorCode, description)
		switch {
		case strings.Contains(description, "can't parse entities"):
			log.Printf(formatted)
			config.Text = text
			config.ParseMode = ""
			goto Retry
		case strings.Contains(description, "Too Many Requests"):
			retryDelay = time.Duration(resp.Parameters.RetryAfter) * time.Second
			functions.Sleep(&retryDelay)
			goto Retry
		case strings.Contains(description, "Forbidden"):
			return 0, e.Forbidden
		case strings.Contains(description, "Bad Request"):
			log.Printf(formatted)
			return 0, err
		case attempts < constants.RetryAttempts:
			functions.Sleep(&retryDelay)
			goto Retry
		default:
			return 0, err
		}
	}
	message := &tgbotapi.Message{}
	err = json.Unmarshal(resp.Result, message)
	if err != nil {
		log.Printf("can't decode result from %d: %s", id(ctx), err)
		return 0, err
	}
	return message.MessageID, nil
}

func Edit(ctx context.Context, text string, messageID int, replyMarkup *tgbotapi.InlineKeyboardMarkup) error {
	formatted := functions.MarkdownV2(text)
	config := tgbotapi.NewEditMessageText(id(ctx), messageID, formatted)
	config.ReplyMarkup = replyMarkup
	config.DisableWebPagePreview = true
	config.ParseMode = tgbotapi.ModeMarkdownV2
	retryDelay := constants.RetryDelay
	attempts := 0
Retry:
	attempts++
	resp, err := bot.Request(config)
	if err != nil {
		description := resp.Description
		log.Printf("can't edit message for %d: %d %s", id(ctx), resp.ErrorCode, description)
		switch {
		case strings.Contains(description, "message to edit not found"):
			return e.UserDeletedMessage
		case strings.Contains(description, "can't parse entities"):
			log.Printf(formatted)
			config.Text = text
			config.ParseMode = ""
			goto Retry
		case strings.Contains(description, "Too Many Requests"):
			retryDelay = time.Duration(resp.Parameters.RetryAfter) * time.Second
			functions.Sleep(&retryDelay)
			goto Retry
		case strings.Contains(description, "Forbidden"):
			return e.Forbidden
		case strings.Contains(description, "Bad Request"):
			log.Printf(formatted)
			return err
		case attempts < constants.RetryAttempts:
			functions.Sleep(&retryDelay)
			goto Retry
		default:
			return err
		}
	}
	return nil
}

func SetTyping(ctx context.Context, isTyping *atomic.Bool) {
	config := tgbotapi.NewChatAction(id(ctx), tgbotapi.ChatTyping)
Loop:
	for isTyping.Load() {
		resp, err := bot.Request(config)
		if err != nil {
			description := resp.Description
			log.Printf("can't set typing for %d: %d %s", id(ctx), resp.ErrorCode, description)
			switch {
			case strings.Contains(description, "Too Many Requests"):
				retryDelay := time.Duration(resp.Parameters.RetryAfter) * time.Second
				functions.Sleep(&retryDelay)
			case strings.Contains(description, "Forbidden"):
				break Loop
			}
		}
		time.Sleep(5800 * time.Millisecond)
	}
}

func Delete(ctx context.Context, messageID int) {
	config := tgbotapi.NewDeleteMessage(id(ctx), messageID)
	resp, err := bot.Request(config)
	if err != nil {
		log.Printf("can't delete message for %d: %d %s", id(ctx), resp.ErrorCode, resp.Description)
	}
}

func SetCommands(ctx context.Context) {
	lang := ctx.Value("language_code").(string)
	scope := tgbotapi.NewBotCommandScopeChat(id(ctx))
	config := tgbotapi.NewSetMyCommandsWithScope(scope, commands[lang]...)
	resp, err := bot.Request(config)
	if err != nil {
		log.Printf("can't set commands for %d: %d %s", id(ctx), resp.ErrorCode, resp.Description)
	}
}

func GetPhotoURL(ctx context.Context, message *tgbotapi.Message) (string, error) {
	photo := message.Photo[len(message.Photo)-1]
	photoURL, err := bot.GetFileDirectURL(photo.FileID)
	if err != nil {
		log.Printf("can't get photo url from %d: %s", id(ctx), err)
		return "", err
	}
	return photoURL, nil
}

func SendPhoto(ctx context.Context, photoURL string) {
	photo := tgbotapi.FileURL(photoURL)
	config := tgbotapi.NewPhoto(id(ctx), photo)
	resp, err := bot.Request(config)
	if err != nil {
		log.Printf("can't send photo to %d: %d %s", id(ctx), resp.ErrorCode, resp.Description)
	}
}
