package telegram

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

func Send(ctx context.Context, text string, replyToMsgID int, addButton bool) (int, error) {
	userID := ctx.Value("user_id").(int64)
	lang := ctx.Value("language_code").(string)

	text = functions.MarkdownV2(text)
	config := tgbotapi.NewMessage(userID, text)
	config.ReplyToMessageID = replyToMsgID
	config.ParseMode = tgbotapi.ModeMarkdownV2
	config.DisableWebPagePreview = true
	if addButton {
		config.ReplyMarkup = newChatButton(lang)
	}

	resp, err := bot.Request(config)
	if err != nil {
		log.Printf("can't send message to %d: %v", userID, err)
		switch {
		case strings.Contains(err.Error(), "bot was blocked by the user"):
			return 0, e.UserBlockedError
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

func Edit(ctx context.Context, text string, messageID int, addButton bool) error {
	userID := ctx.Value("user_id").(int64)
	lang := ctx.Value("language_code").(string)

	text = functions.MarkdownV2(text)
	config := tgbotapi.NewEditMessageText(userID, messageID, text)
	config.ParseMode = tgbotapi.ModeMarkdownV2
	config.DisableWebPagePreview = true
	if addButton {
		config.ReplyMarkup = newChatButton(lang)
	}

	_, err := bot.Request(config)
	if err != nil {
		log.Printf("can't edit message for %d: %v", userID, err)
		switch {
		case strings.Contains(err.Error(), "bot was blocked by the user"):
			return e.UserBlockedError
		case strings.Contains(err.Error(), "message to edit not found"):
			return e.UserDeletedMessage
		}
		return err
	}
	return nil
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
