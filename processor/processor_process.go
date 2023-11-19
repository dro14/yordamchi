package processor

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"time"

	"github.com/dro14/yordamchi/clients/telegram"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/dro14/yordamchi/storage/postgres"
	"github.com/dro14/yordamchi/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) Process(ctx context.Context, message *tgbotapi.Message, isPremium string) {
	utils.SendInfoMessage(message.From.ID, message.MessageID)
	messageID, err := p.telegram.SendMessage(ctx, text.Loading[lang(ctx)], message.MessageID, nil)
	if err != nil {
		log.Println("can't send loading message")
		return
	}

	isTyping := &atomic.Bool{}
	isTyping.Store(true)
	go p.telegram.SetTyping(ctx, isTyping)
	defer isTyping.Store(false)

	msg := &postgres.Message{IsPremium: isPremium}
	msg.Requests++
	msg.FirstSend = int(time.Since(start(ctx)).Milliseconds())

	if message.Photo != nil {
		message.Text, err = p.telegram.PhotoURL(ctx, message)
		if err != nil {
			message.Text = message.Caption
		} else {
			message.Text = message.Text + utils.Delimiter + message.Caption
			msg.IsPremium = "vision"
		}
	}

	msg.Activity = p.redis.IncrementActivity(ctx, message, isPremium)
	defer p.redis.DecrementActivity(ctx)

	i := 0
	var completion string
	var completions []string
	channel := make(chan string)
	if lang(ctx) == "uz" && isPremium == "false" {
		message.Text = p.apis.Translate("auto", "en", message.Text)
		ctx = context.WithValue(ctx, "stream", false)
		go p.openai.ProcessCompletions(ctx, message.Text, msg, channel)
		completion = <-channel
		completion = p.apis.Translate("auto", "uz", completion)
		completions = utils.Slice(completion)

		var replyMarkup *tgbotapi.InlineKeyboardMarkup
		if len(completions) == 1 {
			replyMarkup = p.newChatButton(ctx)
		}

		err = p.telegram.EditMessage(ctx, completions[i], messageID, replyMarkup)
		if err != nil {
			log.Println("can't add new chat button")
		}
		msg.Requests++
		time.Sleep(utils.RequestInterval)

		for i = 1; i < len(completions); i++ {
			if i == len(completions)-1 {
				replyMarkup = p.newChatButton(ctx)
			}
			_, err = p.telegram.SendMessage(ctx, completions[i], 0, replyMarkup)
			if err != nil {
				log.Println("can't send completion")
				i--
			}
			msg.Requests++
			time.Sleep(utils.RequestInterval)
		}
	} else {
		go p.openai.ProcessCompletions(ctx, message.Text, msg, channel)
		for completion = range channel {
			completions = utils.Slice(completion)
			if i >= len(completions) {
				i = len(completions) - 1
			}

			err = p.telegram.EditMessage(ctx, completions[i], messageID, nil)
			if errors.Is(err, telegram.ErrForbidden) {
				return
			} else if errors.Is(err, telegram.ErrDeletedMessage) {
				i--
			}
			msg.Requests++
			time.Sleep(utils.RequestInterval)

			switch completion {
			case text.RequestFailed[lang(ctx)]:
				return
			case text.Error[lang(ctx)]:
				i = 0
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
				msg.Requests++
				time.Sleep(utils.RequestInterval)
			}
		}

		err = p.telegram.EditMessage(ctx, completions[i], messageID, p.newChatButton(ctx))
		if err != nil {
			log.Printf("can't add new chat button")
		}
		msg.Requests++
	}

	p.redis.DecrementRequests(ctx)
	msg.LastEdit = int(time.Since(start(ctx)).Milliseconds())
	msg.CreatedOn = time.Unix(int64(message.Date), 0).Format(time.DateOnly)
	msg.PromptedAt = time.Unix(int64(message.Date), 0).Format(time.TimeOnly)
	msg.CompletedAt = time.Now().Format(time.TimeOnly)
	p.postgres.SaveMessage(ctx, message.From, msg)
	utils.SendInfoMessage(message.From.ID, messageID)
}
