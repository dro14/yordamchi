package processor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/clients/telegram"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/dro14/yordamchi/storage/postgres"
	"github.com/dro14/yordamchi/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) Process(ctx context.Context, message *tgbotapi.Message, isPremium string) {
	messageID, err := p.telegram.SendMessage(ctx, text.Loading[lang(ctx)], message.MessageID, nil)
	if err != nil {
		log.Println("can't send loading message")
		return
	}

	isTyping := &atomic.Bool{}
	isTyping.Store(true)
	go p.telegram.SetTyping(ctx, isTyping)
	defer isTyping.Store(false)

	stats := &postgres.Stats{IsPremium: isPremium}
	stats.Requests++
	beginning := ctx.Value("beginning").(time.Time)
	stats.FirstSend = time.Since(beginning).Milliseconds()
	stats.Activity = p.redis.IncrementActivity(ctx, message, isPremium)
	defer p.redis.DecrementActivity(ctx)

	if message.Photo != nil {
		message.Text, err = p.telegram.PhotoURL(ctx, message)
		if err != nil {
			message.Text = message.Caption
		} else if isPremium == models.GPT4 {
			message.Text = message.Text + "\n\n\n" + message.Caption
			ctx = context.WithValue(ctx, "model", models.GPT4V)
			stats.IsPremium = models.GPT4V
		} else {
			message.Text = p.apis.Vision(ctx, message.Text, message.Caption)
		}
	}

	if message.From.ID == 1792604195 {
		utils.SendInfoMessage(message.Text)
	}

	var tokensUsed int
	var completion string
	var completions []string
	channel := make(chan string)
	if lang(ctx) == "uz" {
		message.Text = p.apis.Translate("auto", "en", message.Text)
		ctx = context.WithValue(ctx, "stream", false)
		go p.openai.ProcessCompletions(ctx, message.Text, stats, channel)
		completion = <-channel
		completion = p.apis.Translate("auto", "uz", completion)

		tokensUsed = stats.PromptTokens + stats.CompletionTokens
		if ctx.Value("model") != models.GPT3 {
			completion = fmt.Sprintf(text.TokensUsed["uz"], completion, tokensUsed)
		}
		completions = utils.Slice(completion)

		var replyMarkup *tgbotapi.InlineKeyboardMarkup
		if len(completions) == 1 {
			replyMarkup = p.newChatButton(lang(ctx))
		}

		err = p.telegram.EditMessage(ctx, completions[0], messageID, replyMarkup)
		if err != nil {
			log.Println("can't add new chat button")
		}
		stats.Requests++

		for i := 1; i < len(completions); i++ {
			if i == len(completions)-1 {
				replyMarkup = p.newChatButton(lang(ctx))
			}
			_, err = p.telegram.SendMessage(ctx, completions[i], 0, replyMarkup)
			if err != nil {
				log.Println("can't send completion")
				i--
			}
			stats.Requests++
		}
	} else {
		go p.openai.ProcessCompletions(ctx, message.Text, stats, channel)

		i := 0
		for completion = range channel {
			completions = utils.Slice(completion)
			if i >= len(completions) {
				i = len(completions) - 1
			}

			err = p.telegram.EditMessage(ctx, completions[i], messageID, nil)
			if errors.Is(err, telegram.ErrForbidden) {
				return
			} else if errors.Is(err, telegram.ErrDeletedMessage) {
				log.Println("user deleted completion")
				i--
			}
			stats.Requests++
			time.Sleep(utils.RequestInterval)

			switch completion {
			case text.RequestFailed[lang(ctx)]:
				return
			case text.Error[lang(ctx)]:
				i--
			}

			for i < len(completions)-1 {
				i++
				messageID, err = p.telegram.SendMessage(ctx, completions[i], 0, nil)
				if errors.Is(err, telegram.ErrForbidden) {
					return
				} else if err != nil {
					log.Println("can't send next message")
					i--
				}
				stats.Requests++
				time.Sleep(utils.RequestInterval)
			}
		}

		tokensUsed = stats.PromptTokens + stats.CompletionTokens
		if ctx.Value("model") != models.GPT3 {
			completions[i] = fmt.Sprintf(text.TokensUsed[lang(ctx)], completions[i], tokensUsed)
		}

		err = p.telegram.EditMessage(ctx, completions[i], messageID, p.newChatButton(lang(ctx)))
		if err != nil {
			log.Printf("can't add new chat button")
		}
		stats.Requests++
	}

	stats.LastEdit = time.Since(beginning).Milliseconds()
	stats.CompletedAt = time.Now().Unix()
	p.redis.Decrement(ctx, tokensUsed)
	p.postgres.SaveMessage(ctx, message.From, stats)
	if message.From.ID == 1792604195 {
		utils.SendInfoMessage(completion)
	}
}
