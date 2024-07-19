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

func (t *Telegram) SendMessage(ctx context.Context, text string, replyToMsgID int, replyMarkup any) (int, error) {
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

func (t *Telegram) SetTyping(ctx context.Context) *atomic.Bool {
	isTyping := &atomic.Bool{}
	isTyping.Store(true)

	go func() {
		config := tgbotapi.NewChatAction(id(ctx), tgbotapi.ChatTyping)
		for isTyping.Load() {
			_, err := t.makeRequest(ctx, config)
			if errors.Is(err, ErrForbidden) {
				break
			}
			time.Sleep(5800 * time.Millisecond)
		}
	}()

	return isTyping
}

func (t *Telegram) DeleteMessage(ctx context.Context, messageID int) {
	config := tgbotapi.NewDeleteMessage(id(ctx), messageID)
	_, _ = t.makeRequest(ctx, config)
}

func (t *Telegram) SetCommands(ctx context.Context) {
	scope := tgbotapi.NewBotCommandScopeChat(id(ctx))
	config := tgbotapi.NewSetMyCommandsWithScope(scope, commands[lang(ctx)]...)
	_, err := t.makeRequest(ctx, config)
	if err != nil {
		log.Printf("user %d: can't set commands", id(ctx))
	}
}

func (t *Telegram) PhotoURL(ctx context.Context, photos []tgbotapi.PhotoSize) string {
	config := tgbotapi.FileConfig{FileID: photos[len(photos)-1].FileID}
	resp, err := t.makeRequest(ctx, config)
	if err != nil {
		log.Printf("user %d: can't get photo url", id(ctx))
		return ""
	}

	file := &tgbotapi.File{}
	err = json.Unmarshal(resp.Result, file)
	if err != nil {
		log.Printf("user %d: can't decode result: %s", id(ctx), err)
		return ""
	}
	return file.Link(t.token)
}

func (t *Telegram) SendPhoto(ctx context.Context, path, caption string, replyMarkup *tgbotapi.InlineKeyboardMarkup) error {
	caption = utils.Slice(caption, 1024)[0]
	config := tgbotapi.NewPhoto(id(ctx), tgbotapi.FilePath(path))
	config.Caption = utils.MarkdownV2(caption)
	config.ReplyMarkup = replyMarkup
	config.ParseMode = tgbotapi.ModeMarkdownV2
Retry:
	_, err := t.makeRequest(ctx, config)
	if errors.Is(err, ErrMarkdown) {
		log.Println(config.Caption)
		config.Caption = caption
		config.ParseMode = ""
		goto Retry
	} else if err != nil {
		return err
	}
	return nil
}

func (t *Telegram) AnswerCallbackQuery(ctx context.Context, ID, text string) {
	config := tgbotapi.NewCallback(ID, text)
	_, _ = t.makeRequest(ctx, config)
}

func (t *Telegram) SetKeyboard(ctx context.Context, text string, questions []string) error {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(questions[0])),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(questions[1])),
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(questions[2])))
	keyboard.OneTimeKeyboard = true

	formatted := utils.MarkdownV2(text)
	config := tgbotapi.NewMessage(id(ctx), formatted)
	config.ReplyMarkup = keyboard
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
