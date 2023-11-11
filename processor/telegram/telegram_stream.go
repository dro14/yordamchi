package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dro14/yordamchi/client/ocr"
	"github.com/dro14/yordamchi/client/telegram"
	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/lib/models"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/postgres"
	"github.com/dro14/yordamchi/processor/openai"
	"github.com/dro14/yordamchi/processor/telegram/button"
	"github.com/dro14/yordamchi/processor/telegram/info_bot"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Process(ctx context.Context, message *tgbotapi.Message, isPremium string) {

	if message.From.ID == 1792604195 {
		info_bot.Send(message.Text)
	}

	if strings.Contains(message.Text, "#image") {
		GeneratePhoto(ctx, message)
		return
	}

	stats := &types.Stats{IsPremium: isPremium}

	stats.Requests++
	messageID, err := telegram.Send(ctx, text.Loading[lang(ctx)], message.MessageID, nil)
	if err != nil {
		log.Printf("can't send loading message")
		return
	}

	beginning := ctx.Value("beginning").(time.Time)
	stats.FirstSend = time.Since(beginning).Milliseconds()

	stats.Activity = redis.IncrementActivity(ctx, message, isPremium)
	defer redis.DecrementActivity(ctx)

	isTyping := &atomic.Bool{}
	isTyping.Store(true)
	go telegram.SetTyping(ctx, isTyping)
	defer isTyping.Store(false)

	if message.Photo != nil {
		photoURL, err := telegram.GetPhotoURL(message)
		if err != nil {
			log.Printf("can't get photo url")
			message.Text = message.Caption
		} else if isPremium == models.GPT4 {
			ctx = context.WithValue(ctx, "model", models.GPT4V)
			message.Text = photoURL + "\n\n\n" + message.Caption
		} else {
			message.Text = ocr.Read(ctx, photoURL, message.Caption)
		}
	}

	if lang(ctx) == "uz" {
		completions := UseTranslator(ctx, message, stats)

		var replyMarkup *tgbotapi.InlineKeyboardMarkup
		if len(completions) == 1 {
			replyMarkup = button.NewChat(lang(ctx))
		}

		stats.Requests++
		err = telegram.Edit(ctx, completions[0], messageID, replyMarkup)
		if err != nil {
			log.Printf("can't add new chat button")
		}

		for i := 1; i < len(completions); i++ {
			if i == len(completions)-1 {
				replyMarkup = button.NewChat(lang(ctx))
			}
			stats.Requests++
			_, err = telegram.Send(ctx, completions[i], 0, replyMarkup)
			if err != nil {
				log.Printf("can't send completion")
				i--
			}
		}

		stats.LastEdit = time.Since(beginning).Milliseconds()
		stats.CompletedAt = time.Now().Unix()
		postgres.SaveMessage(ctx, stats, message.From)
	} else {
		messages := redis.LoadContext(ctx, message.Text)
		channel := make(chan string)
		go openai.ProcessWithStream(ctx, messages, stats, channel)

		index := 0
		completion := ""
		var completions []string
		for completion = range channel {
			completions = functions.Slice(completion)
			if index >= len(completions) {
				index = len(completions) - 1
			}

			stats.Requests++
			err = telegram.Edit(ctx, completions[index], messageID, nil)
			if errors.Is(err, e.UserBlockedBot) {
				return
			} else if errors.Is(err, e.UserDeletedMessage) {
				log.Printf("user deleted completion")
				index--
			}

			switch completion {
			case text.TooLong[lang(ctx)]:
				log.Printf("prompt was too long")
				return
			case text.RequestFailed[lang(ctx)]:
				return
			case text.Error[lang(ctx)]:
				index--
			}

			for index < len(completions)-1 {
				index = len(completions) - 1
				stats.Requests++
				time.Sleep(constants.RequestInterval)
				messageID, err = telegram.Send(ctx, completions[index], 0, nil)
				if errors.Is(err, e.UserBlockedBot) {
					return
				} else if err != nil {
					log.Printf("can't send next message")
					index--
				}
			}

			time.Sleep(constants.RequestInterval)
		}

		tokensUsed := stats.PromptTokens + stats.CompletionTokens
		if ctx.Value("model") == models.GPT4 {
			completions[index] = fmt.Sprintf(text.TokensUsed[lang(ctx)], completions[index], tokensUsed)
		}

		stats.Requests++
		err = telegram.Edit(ctx, completions[index], messageID, button.NewChat(lang(ctx)))
		if err != nil {
			log.Printf("can't add new chat button")
		}
		stats.LastEdit = time.Since(beginning).Milliseconds()
		stats.CompletedAt = time.Now().Unix()

		redis.Decrement(ctx, tokensUsed)
		redis.StoreContext(ctx, message.Text, completion)
		postgres.SaveMessage(ctx, stats, message.From)

		if message.From.ID == 1792604195 {
			info_bot.Send(completion)
		}
	}
}
