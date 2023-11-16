package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync/atomic"
	"time"

	"github.com/dro14/yordamchi/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Telegram) SendMessage(ctx context.Context, text string, replyToMsgID int, replyMarkup *tgbotapi.InlineKeyboardMarkup) (int, error) {
	formatted := utils.MarkdownV2(text)
	config := tgbotapi.NewMessage(id(ctx), formatted)
	config.ReplyToMessageID = replyToMsgID
	config.ReplyMarkup = replyMarkup
	config.DisableWebPagePreview = true
	config.ParseMode = tgbotapi.ModeMarkdownV2
Retry:
	resp, err := t.makeRequest(ctx, config)
	if errors.Is(err, ErrMarkdown) {
		log.Println(formatted)
		config.Text = text
		config.ParseMode = ""
		goto Retry
	} else if err != nil {
		return 0, err
	}

	message := &tgbotapi.Message{}
	err = json.Unmarshal(resp.Result, message)
	if err != nil {
		log.Printf("user %d: can't decode result: %s", id(ctx), err)
		return 0, err
	}
	return message.MessageID, nil
}

func (t *Telegram) EditMessage(ctx context.Context, text string, messageID int, replyMarkup *tgbotapi.InlineKeyboardMarkup) error {
	formatted := utils.MarkdownV2(text)
	config := tgbotapi.NewEditMessageText(id(ctx), messageID, formatted)
	config.ReplyMarkup = replyMarkup
	config.DisableWebPagePreview = true
	config.ParseMode = tgbotapi.ModeMarkdownV2
Retry:
	_, err := t.makeRequest(ctx, config)
	if errors.Is(err, ErrMarkdown) {
		log.Println(formatted)
		config.Text = text
		config.ParseMode = ""
		goto Retry
	} else if err != nil {
		return err
	}
	return nil
}

func (t *Telegram) SetTyping(ctx context.Context, isTyping *atomic.Bool) {
	config := tgbotapi.NewChatAction(id(ctx), tgbotapi.ChatTyping)
	for isTyping.Load() {
		_, err := t.makeRequest(ctx, config)
		if errors.Is(err, ErrForbidden) {
			break
		}
		time.Sleep(5800 * time.Millisecond)
	}
}

func (t *Telegram) DeleteMessage(ctx context.Context, messageID int) {
	config := tgbotapi.NewDeleteMessage(id(ctx), messageID)
	_, err := t.makeRequest(ctx, config)
	if err != nil {
		log.Printf("user %d: can't delete message", id(ctx))
	}
}

func (t *Telegram) SetCommands(ctx context.Context) {
	lang := ctx.Value("language_code").(string)
	scope := tgbotapi.NewBotCommandScopeChat(id(ctx))
	config := tgbotapi.NewSetMyCommandsWithScope(scope, commands[lang]...)
	_, err := t.makeRequest(ctx, config)
	if err != nil {
		log.Printf("user %d: can't set commands", id(ctx))
	}
}

func (t *Telegram) PhotoURL(ctx context.Context, message *tgbotapi.Message) (string, error) {
	photo := message.Photo[len(message.Photo)-1]
	config := tgbotapi.FileConfig{FileID: photo.FileID}
	resp, err := t.makeRequest(ctx, config)
	if err != nil {
		log.Printf("user %d: can't get photo url", id(ctx))
		return "", err
	}

	file := &tgbotapi.File{}
	err = json.Unmarshal(resp.Result, file)
	if err != nil {
		log.Printf("user %d: can't decode result: %s", id(ctx), err)
		return "", err
	}
	return file.Link(t.token), nil
}

func (t *Telegram) SendPhoto(ctx context.Context, photoURL string) {
	photo := tgbotapi.FileURL(photoURL)
	config := tgbotapi.NewPhoto(id(ctx), photo)
	_, err := t.makeRequest(ctx, config)
	if err != nil {
		log.Printf("user %d: can't send photo", id(ctx))
	}
}

func (t *Telegram) SendFile(ctx context.Context, filepath string) {
	file := tgbotapi.FilePath(filepath)
	config := tgbotapi.NewDocument(id(ctx), file)
	_, err := t.makeRequest(ctx, config)
	if err != nil {
		log.Printf("user %d: can't send file", id(ctx))
	}
}
