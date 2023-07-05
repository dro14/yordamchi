package telegram

import (
	"context"
	"fmt"
	"log"

	"github.com/dro14/yordamchi/client/telegram"
	"github.com/dro14/yordamchi/postgres"
	"github.com/dro14/yordamchi/processor/telegram/button"
	"github.com/dro14/yordamchi/processor/telegram/info_bot"
	"github.com/dro14/yordamchi/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func blocked(ctx context.Context) {

	_, err := telegram.SendMessage(ctx, text.Blocked[lang(ctx)], 0, button.Blocked(lang(ctx)))
	if err != nil {
		log.Printf("can't send blocked message")
	}
}

func deactivated(ctx context.Context, user *tgbotapi.User) {

	postgres.DeactivateUser(ctx, user)
	joinedAt := format(postgres.JoinedAt(ctx, user))
	deactivatedAt := format(postgres.DeactivatedAt(ctx, user))
	rejoinedAt := format(postgres.RejoinedAt(ctx, user))

	if len(user.UserName) > 0 {
		user.UserName = "@" + user.UserName
	}

	info_bot.Send(fmt.Sprintf(
		`is_active: 🚫

id: %d
first_name: %s
last_name: %s
username: %s
language_code: %s

joined_at:           %s
deactivated_at: %s
rejoined_at:       %s`, user.ID, user.FirstName, user.LastName, user.UserName, user.LanguageCode, joinedAt, deactivatedAt, rejoinedAt))
}

func rejoined(ctx context.Context, user *tgbotapi.User) {

	postgres.RejoinUser(ctx, user)
	joinedAt := format(postgres.JoinedAt(ctx, user))
	deactivatedAt := format(postgres.DeactivatedAt(ctx, user))
	rejoinedAt := format(postgres.RejoinedAt(ctx, user))

	if len(user.UserName) > 0 {
		user.UserName = "@" + user.UserName
	}

	info_bot.Send(fmt.Sprintf(
		`is_active: ✅

id: %d
first_name: %s
last_name: %s
username: %s
language_code: %s

joined_at:           %s
deactivated_at: %s
rejoined_at:       %s`, user.ID, user.FirstName, user.LastName, user.UserName, user.LanguageCode, joinedAt, deactivatedAt, rejoinedAt))
}
