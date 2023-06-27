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
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/postgres"
	"github.com/dro14/yordamchi/processor/openai"
	"github.com/dro14/yordamchi/processor/telegram/button"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Stream(ctx context.Context, message *tgbotapi.Message, isPremium string) {

	messages := redis.LoadContext(ctx, message.Text)
	stats := &types.Stats{IsPremium: isPremium}
	channel := make(chan string)
	go openai.Process(ctx, messages, stats, channel)

	stats.Activity = redis.IncrementActivity(ctx, message, isPremium)
	defer redis.DecrementActivity(ctx)

	stats.Requests++
	messageID, err := telegram.SendMessage(ctx, text.Loading[lang(ctx)], message.MessageID, nil)
	if err != nil {
		log.Printf("can't send loading message")
		return
	}
	beginning := ctx.Value("beginning").(time.Time)
	stats.FirstSend = time.Since(beginning).Milliseconds()

	isTyping := &atomic.Bool{}
	isTyping.Store(true)
	go telegram.SetTyping(ctx, isTyping)
	defer isTyping.Store(false)

	index := 0
	completion := ""
	var completions []string
	for completion = range channel {

		completions = slice(completion)
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
			fallthrough
		case text.RequestFailed[lang(ctx)]:
			return
		case text.Error[lang(ctx)]:
			index--
		}

		if index < len(completions)-1 {
			time.Sleep(constants.RequestInterval)
			index++
			stats.Requests++
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
	err = telegram.EditMessage(ctx, completions[index], messageID, button.NewChat(lang(ctx)))
	if err != nil {
		log.Printf("can't add new chat button")
	}
	stats.LastEdit = time.Since(beginning).Milliseconds()
	stats.CompletedAt = time.Now().Unix()

	redis.Decrement(ctx, tokensUsed)
	redis.StoreContext(ctx, message.Text, completion)
	postgres.SaveMessage(ctx, stats, message.From)
}
