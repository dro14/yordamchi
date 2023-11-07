package telegram

import (
	"context"
	"log"

	"github.com/dro14/yordamchi/lib/functions"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

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
		log.Printf("can't edit message for %d: %v\nmessage: %q", userID, err, message)
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

func SetCommands(ctx context.Context, lang string) {
	userID := ctx.Value("user_id").(int64)
	scope := tgbotapi.NewBotCommandScopeChat(userID)
	config := tgbotapi.NewSetMyCommandsWithScope(scope, commands[lang]...)
	_, err := bot.Request(config)
	if err != nil {
		log.Printf("can't set commands for %d: %v", userID, err)
	}
}
