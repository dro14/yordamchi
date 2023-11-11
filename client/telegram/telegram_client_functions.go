package telegram

import (
	"context"
	"encoding/json"
	"log"
	"sync/atomic"
	"time"

	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Send(ctx context.Context, text string, replyToMsgID int, replyMarkup *tgbotapi.InlineKeyboardMarkup) (int, error) {
	userID := ctx.Value("user_id").(int64)
	text = functions.MarkdownV2(text)
	config := tgbotapi.NewMessage(userID, text)
	config.ReplyToMessageID = replyToMsgID
	config.ReplyMarkup = replyMarkup
	config.DisableWebPagePreview = true
	config.ParseMode = tgbotapi.ModeMarkdownV2

	resp, err := bot.Request(config)
	if err != nil {
		log.Printf("can't send message to %d: %v", userID, err)
		switch err.Error() {
		case "Forbidden: bot was blocked by the user":
			return 0, e.UserBlockedBot
		}
		return 0, err
	}

	message := &tgbotapi.Message{}
	err = json.Unmarshal(resp.Result, message)
	if err != nil {
		log.Printf("can't decode message for %d: %v", userID, err)
	}
	return message.MessageID, nil
}

func Edit(ctx context.Context, text string, messageID int, replyMarkup *tgbotapi.InlineKeyboardMarkup) error {
	userID := ctx.Value("user_id").(int64)
	text = functions.MarkdownV2(text)
	config := tgbotapi.NewEditMessageText(userID, messageID, text)
	config.ReplyMarkup = replyMarkup
	config.DisableWebPagePreview = true
	config.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := bot.Request(config)
	if err != nil {
		log.Printf("can't edit message for %d: %v", userID, err)
		switch err.Error() {
		case "Forbidden: bot was blocked by the user":
			return e.UserBlockedBot
		case "Bad Request: message to edit not found":
			return e.UserDeletedMessage
		}
		return err
	}
	return nil
}

func SetTyping(ctx context.Context, isTyping *atomic.Bool) {
	userID := ctx.Value("user_id").(int64)
	config := tgbotapi.NewChatAction(userID, tgbotapi.ChatTyping)
Loop:
	for isTyping.Load() {
		_, err := bot.Request(config)
		if err != nil {
			log.Printf("can't set typing for %d: %v", userID, err)
			switch err.Error() {
			case "Forbidden: bot was blocked by the user":
				break Loop
			}
		}

		time.Sleep(5800 * time.Millisecond)
	}
}

func Delete(ctx context.Context, messageID int) error {
	userID := ctx.Value("user_id").(int64)
	config := tgbotapi.NewDeleteMessage(userID, messageID)
	_, err := bot.Request(config)
	if err != nil {
		log.Printf("can't delete message for %d: %v", userID, err)
		return err
	}
	return nil
}

func SetCommands(ctx context.Context) {
	userID := ctx.Value("user_id").(int64)
	lang := ctx.Value("language_code").(string)
	scope := tgbotapi.NewBotCommandScopeChat(userID)
	config := tgbotapi.NewSetMyCommandsWithScope(scope, commands[lang]...)
	_, err := bot.Request(config)
	if err != nil {
		log.Printf("can't set commands for %d: %v", userID, err)
	}
}

func GetPhotoURL(message *tgbotapi.Message) (string, error) {
	photo := message.Photo[len(message.Photo)-1]
	photoURL, err := bot.GetFileDirectURL(photo.FileID)
	if err != nil {
		log.Printf("can't get photo url: %v", err)
		return "", err
	}
	return photoURL, nil
}
