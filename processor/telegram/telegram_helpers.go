package telegram

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
	"github.com/gotd/td/tg"
)

var blockedUsers sync.Map

func isBlocked(userID int64) bool {
	_, ok := blockedUsers.Load(userID)
	if ok {
		return true
	} else {
		blockedUsers.Store(userID, true)
		return false
	}
}

func messageUpdate(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) (context.Context, *tg.Message, *tg.User) {

	start := time.Now()
	ctx = context.WithValue(ctx, "start", start)

	message, ok := update.Message.(*tg.Message)
	if !ok || message.Out || len(strings.TrimSpace(message.Message)) == 0 {
		return ctx, nil, nil
	}

	peerUser, ok := message.PeerID.(*tg.PeerUser)
	if !ok {
		return ctx, nil, nil
	}

	user := entities.Users[peerUser.UserID]
	if user.Bot || isBlocked(user.ID) {
		return ctx, nil, nil
	}

	ctx = context.WithValue(ctx, "date", message.Date)
	ctx = context.WithValue(ctx, "user_id", user.ID)
	ctx = context.WithValue(ctx, "language_code", functions.LanguageCode(user.LangCode))
	return ctx, message, user
}

func callbackUpdate(ctx context.Context, entities tg.Entities, update *tg.UpdateBotCallbackQuery) (context.Context, string) {
	user := entities.Users[update.UserID]
	ctx = context.WithValue(ctx, "message_id", update.MsgID)
	ctx = context.WithValue(ctx, "user_id", user.ID)
	ctx = context.WithValue(ctx, "language_code", functions.LanguageCode(user.LangCode))
	return ctx, string(update.Data)
}

func botStoppedUpdate(ctx context.Context, entities tg.Entities, update *tg.UpdateBotStopped) (context.Context, *tg.User) {
	user := entities.Users[update.UserID]
	ctx = context.WithValue(ctx, "date", update.Date)
	ctx = context.WithValue(ctx, "user_id", user.ID)
	ctx = context.WithValue(ctx, "language_code", functions.LanguageCode(user.LangCode))
	return ctx, user
}

func command(message *tg.Message) string {

	entities := message.Entities
	if len(entities) == 0 {
		return ""
	}

	entity, ok := entities[0].(*tg.MessageEntityBotCommand)
	if !ok || entity.Offset != 0 || entity.Length == 0 {
		return ""
	}

	botCommand := message.Message[1:entity.Length]
	if i := strings.Index(botCommand, "@"); i != -1 {
		botCommand = botCommand[:i]
	}

	return botCommand
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

func format(timestamp string) string {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return timestamp
	}
	return t.Format("2006-01-02 15:04:05")
}

func slice(completion string) []string {

	var completions []string

	for len(completion) > 4096 {
		cutIndex := 0
		for i := 4096; i >= 0; i-- {
			if completion[i] == ' ' || completion[i] == '\n' || completion[i] == '\t' || completion[i] == '\r' {
				cutIndex = i
				break
			}
		}
		completions = append(completions, completion[:cutIndex])
		completion = completion[cutIndex:]
	}

	return append(completions, completion)
}

func msg(ctx context.Context, lang string) string {

	if redis.UserStatus(ctx) == types.PremiumStatus {
		return fmt.Sprintf(text.Settings[lang], text.PremiumTariff[lang], text.Unlimited[lang], redis.Expiration(ctx))
	}

	return fmt.Sprintf(text.Settings[lang], text.FreeTariff[lang], redis.Requests(ctx), redis.Expiration(ctx))
}
