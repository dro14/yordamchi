package telegram

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/dro14/yordamchi/client/telegram"
	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/ocr"
	"github.com/dro14/yordamchi/postgres"
	"github.com/dro14/yordamchi/processor/openai"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Stream(ctx context.Context, message *tgbotapi.Message, isPremium string) {

	stats := &types.Stats{IsPremium: isPremium}
	stats.Requests++
	messageID, err := telegram.SendMessage(ctx, text.Loading[lang(ctx)], message.MessageID, nil)
	if err != nil {
		log.Printf("can't send loading message")
		return
	}
	beginning := ctx.Value("beginning").(time.Time)
	stats.FirstSend = time.Since(beginning).Milliseconds()
	stats.Activity = redis.IncrementActivity(ctx, message, isPremium)
	defer redis.DecrementActivity(ctx)

	if message.Photo != nil {
		message.Text = ocr.Analyze(ctx, message)
	}

	isTyping := &atomic.Bool{}
	isTyping.Store(true)
	go telegram.SetTyping(ctx, isTyping)
	defer isTyping.Store(false)

	if ctx.Value("target_lang") != "-" {
		completions := UseTranslator(ctx, message, stats)

		stats.Requests++
		err = telegram.Edit(ctx, completions[0], messageID, len(completions) == 1)
		if err != nil {
			log.Printf("can't add new chat button")
		}

		for i := 1; i < len(completions); i++ {
			stats.Requests++
			err = telegram.Send(ctx, completions[i], i == len(completions)-1)
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
			err = telegram.EditMessage(ctx, completions[index], messageID, nil)
			if err == e.UserBlockedError {
				return
			} else if err == e.UserDeletedMessage {
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
				messageID, err = telegram.SendMessage(ctx, completions[index], 0, nil)
				if err == e.UserBlockedError {
					return
				} else if err != nil {
					log.Printf("can't send next message")
					index--
				}
			}

			time.Sleep(constants.RequestInterval)
		}

		tokensUsed := stats.PromptTokens + stats.CompletionTokens
		if ctx.Value("model") == "gpt-4" {
			completions[index] = fmt.Sprintf(text.TokensUsed[lang(ctx)], completions[index], tokensUsed)
		}

		stats.Requests++
		err = telegram.Edit(ctx, completions[index], messageID, true)
		if err != nil {
			log.Printf("can't add new chat button")
		}
		stats.LastEdit = time.Since(beginning).Milliseconds()
		stats.CompletedAt = time.Now().Unix()

		redis.Decrement(ctx, tokensUsed)
		redis.StoreContext(ctx, message.Text, completion)
		postgres.SaveMessage(ctx, stats, message.From)
	}
}
