package processor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dro14/yordamchi/processor/text"
	"github.com/dro14/yordamchi/utils"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) exhausted(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.Exhausted[lang(ctx)], 0, p.settingsButton(ctx))
	if err != nil {
		log.Println("can't send exhausted message")
	}
}

func (p *Processor) processFile(ctx context.Context, message *tgbotapi.Message) {
	pieces := strings.Split(message.Document.FileName, ".")
	switch pieces[len(pieces)-1] {
	case "png", "jpg", "jpeg":
		message.Photo = []tgbotapi.PhotoSize{{FileID: message.Document.FileID}}
		p.process(ctx, message)
		return
	}

	messageID, err := p.telegram.SendMessage(ctx, text.Loading[lang(ctx)], message.MessageID, nil)
	if err != nil {
		log.Println("can't send loading message")
		return
	}

	isTyping := p.telegram.SetTyping(ctx)
	defer isTyping.Store(false)

	var Text string
	errMsg := p.service.Load(ctx, message.Document)
	if supported, found := strings.CutPrefix(errMsg, "supported file formats:"); found {
		Text = fmt.Sprintf(text.UnsupportedFormat[lang(ctx)], supported)
	} else if errMsg != "" {
		Text = text.FailedRequest[lang(ctx)]
	} else {
		Text = fmt.Sprintf(text.FileLoaded[lang(ctx)], message.Document.FileName)
	}

	err = p.telegram.EditMessage(ctx, Text, messageID, nil)
	if err != nil {
		log.Println("can't edit process file message")
	}
}

func (p *Processor) notify(ctx context.Context) {
	patterns := []string{"lang:*", "context:*", "unlimited:*", "premium:*"}
	for {
		notificationTime := p.redis.NotificationTime(ctx)
		time.Sleep(time.Until(notificationTime))
		notificationTime = notificationTime.Add(utils.NotificationInterval)
		p.redis.SetNotificationTime(ctx, notificationTime)

		for _, pattern := range patterns {
			for _, userID := range p.redis.SoonExpires(ctx, pattern) {
				ctx = context.WithValue(ctx, "user_id", userID)
				ctx, _ = p.redis.Lang(ctx, "uz")
				name := p.postgres.User(ctx, &tgbotapi.User{ID: userID})

				var Text string
				var replyMarkup *tgbotapi.InlineKeyboardMarkup
				switch pattern {
				case "lang:*", "context:*":
					if name == "" {
						name = text.DearUser[lang(ctx)]
					}
					Text = fmt.Sprintf(text.Notification1[lang(ctx)], name)
					replyMarkup = p.startButton(ctx)
				case "unlimited:*", "premium:*":
					if name == "" {
						name = text.User[lang(ctx)]
					}
					var subscription string
					if pattern == "unlimited:*" {
						subscription = text.UnlimitedSubscription[lang(ctx)]
					} else {
						subscription = text.PremiumSubscription[lang(ctx)]
					}
					Text = fmt.Sprintf(text.Notification2[lang(ctx)], name, subscription, p.redis.Expiration(ctx))
					replyMarkup = p.settingsButton(ctx)
				}

				_, _ = p.telegram.SendMessage(ctx, Text, 0, replyMarkup)
			}
		}
	}
}
