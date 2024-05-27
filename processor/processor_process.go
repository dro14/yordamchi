package processor

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/clients/telegram"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/dro14/yordamchi/storage/postgres"
	"github.com/dro14/yordamchi/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) process(ctx context.Context, message *tgbotapi.Message, Type string) {
	messageID, err := p.telegram.SendMessage(ctx, text.Loading[lang(ctx)], message.MessageID, nil)
	if err != nil {
		log.Println("can't send loading message")
		return
	}

	isTyping := p.telegram.SetTyping(ctx)
	defer isTyping.Store(false)

	msg := &postgres.Message{Type: Type}
	msg.Requests++
	msg.Input = message.Text
	msg.FirstSend = int(time.Since(start(ctx)).Milliseconds())
	msg.Activity = int(p.activity.Add(1))
	defer p.activity.Add(-1)

	if message.Photo != nil {
		msg.Input = message.Caption
		var photoURL string
		photoURL, err = p.telegram.PhotoURL(ctx, message.Photo)
		if err != nil {
			message.Text = message.Caption
		} else if model(ctx) == models.GPT3 {
			message.Text = p.apis.OCR(ctx, photoURL, message.Caption)
			msg.Type = "ocr"
			msg.Input = message.Text
		} else {
			message.Text = photoURL + utils.Delim + message.Caption
			msg.Type = "vision"
		}
		if message.From.ID == 1792604195 {
			path := message.Photo[len(message.Photo)-1].FileID + ".jpeg"
			_ = utils.DownloadFile(photoURL, path)
			utils.SendInfoMessage(message.Caption, path)
		}
	} else {
		if message.From.ID == 1792604195 {
			utils.SendInfoMessage(message.Text, "")
		}
	}

	i := 0
	var completion string
	var completions []string
	channel := make(chan string)
	if model(ctx) == models.GPT3 && lang(ctx) == "uz" {
		ctx = context.WithValue(ctx, "stream", false)
		go p.openai.ProcessCompletions(ctx, message.Text, msg, channel)
		completion = <-channel
		completion = utils.LaTeX(completion)
		completion = p.apis.Translate("en", "uz", completion)
		completions = utils.Slice(completion, 4096)

		var replyMarkup *tgbotapi.InlineKeyboardMarkup
		if len(completions) == 1 {
			replyMarkup = p.newChatButton(ctx)
		}

		err = p.telegram.EditMessage(ctx, completions[i], messageID, replyMarkup)
		if errors.Is(err, telegram.ErrForbidden) {
			return
		} else if errors.Is(err, telegram.ErrNotFound) {
			i--
		}
		msg.Requests++
		time.Sleep(utils.ReqInterval)

		for i++; i < len(completions); i++ {
			if i == len(completions)-1 {
				replyMarkup = p.newChatButton(ctx)
			}
			_, err = p.telegram.SendMessage(ctx, completions[i], 0, replyMarkup)
			if errors.Is(err, telegram.ErrForbidden) {
				return
			} else if err != nil {
				log.Println("can't send completion")
				i--
			}
			msg.Requests++
			time.Sleep(utils.ReqInterval)
		}
	} else {
		go p.openai.ProcessCompletions(ctx, message.Text, msg, channel)
		for completion = range channel {
			completions = utils.Slice(completion, 4096)
			if i >= len(completions) {
				i = len(completions) - 1
			}

			err = p.telegram.EditMessage(ctx, completions[i], messageID, nil)
			if errors.Is(err, telegram.ErrForbidden) {
				return
			} else if errors.Is(err, telegram.ErrNotFound) {
				i--
			}
			msg.Requests++
			time.Sleep(utils.ReqInterval)

			switch completion {
			case text.BadRequest[lang(ctx)], text.RequestFailed[lang(ctx)]:
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
				time.Sleep(utils.ReqInterval)
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
	if message.From.ID == 1792604195 {
		utils.SendInfoMessage(completion, "")
	}
}
