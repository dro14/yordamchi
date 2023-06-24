package telegram

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/postgres"
	"github.com/dro14/yordamchi/processor/telegram/button"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
	"github.com/gotd/td/tg"
)

func (p *Processor) Stream(ctx context.Context, message *tg.Message, user *tg.User, isPremium bool) {

	messages := redis.LoadContext(ctx, message.Message)
	stats := &types.Stats{IsPremium: isPremium}
	channel := make(chan string)
	go p.Processor.Process(ctx, messages, stats, channel)

	stats.Activity = redis.IncrementActivity(ctx, message, user, isPremium)
	defer redis.DecrementActivity(ctx)

	stats.Requests++
	messageID, err := p.Client.SendMessage(ctx, text.Loading[lang(ctx)], message.ID, nil)
	if err != nil {
		log.Printf("can't send loading message")
		return
	}
	start := ctx.Value("start").(time.Time)
	stats.FirstSend = time.Since(start).Milliseconds()

	isTyping := &atomic.Bool{}
	isTyping.Store(true)
	go p.Client.SetTyping(ctx, isTyping)
	defer isTyping.Store(false)

	index := 0
	completion := ""
	for completion = range channel {

		completions := slice(completion)
		if index >= len(completions) {
			index = len(completions) - 1
		}

		stats.Requests++
		err = p.Client.EditMessage(ctx, completions[index], messageID, nil)
		if err == e.UserBlockedError {
			log.Printf("user blocked bot")
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
			log.Printf("request failed")
			return
		case text.Error[lang(ctx)]:
			log.Printf("stream error")
			index--
		}

		if index < len(completions)-1 {
			time.Sleep(constants.RequestInterval)
			index++
			stats.Requests++
			messageID, err = p.Client.SendMessage(ctx, completions[index], 0, nil)
			if err == e.UserBlockedError {
				log.Printf("user blocked bot")
				return
			} else if err != nil {
				log.Printf("can't send next message")
				index--
			}
		}

		time.Sleep(constants.RequestInterval)
	}

	stats.Requests++
	err = p.Client.EditMessage(ctx, "", messageID, button.NewChat(lang(ctx)))
	if err != nil {
		log.Printf("can't add new chat button")
	}
	stats.LastEdit = time.Since(start).Milliseconds()
	stats.CompletedAt = time.Now().Unix()

	redis.Decrement(ctx)
	redis.StoreContext(ctx, message.Message, completion)
	postgres.SaveMessage(ctx, stats, user)
}
