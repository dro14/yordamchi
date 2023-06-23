package telegram

import (
	"context"
	"fmt"
	"log"

	"github.com/dro14/yordamchi/postgres"
	"github.com/dro14/yordamchi/processor/telegram/button"
	"github.com/dro14/yordamchi/processor/telegram/info"
	"github.com/dro14/yordamchi/text"
	"github.com/gotd/td/tg"
)

func (p *Processor) blocked(ctx context.Context) {

	_, err := p.Client.SendMessage(ctx, text.Blocked[lang(ctx)], 0, button.Blocked(lang(ctx)))
	if err != nil {
		log.Printf("can't send blocked message")
	}
}

func (p *Processor) exhausted(ctx context.Context) {

	_, err := p.Client.SendMessage(ctx, text.Exhausted[lang(ctx)], 0, button.Exhausted(lang(ctx)))
	if err != nil {
		log.Printf("can't send exhausted message")
	}
}

func (p *Processor) deactivated(ctx context.Context, user *tg.User) {

	postgres.DeactivateUser(ctx, user)
	joinedAt := format(postgres.JoinedAt(ctx, user))
	deactivatedAt := format(postgres.DeactivatedAt(ctx, user))
	rejoinedAt := format(postgres.RejoinedAt(ctx, user))

	if len(user.Username) > 0 {
		user.Username = "@" + user.Username
	}

	info.SendInfoMessage(ctx, fmt.Sprintf(
		`is_active: 🚫

id: %d
first_name: %s
last_name: %s
username: %s
language_code: %s

joined_at:           %s
deactivated_at: %s
rejoined_at:       %s`, user.ID, user.FirstName, user.LastName, user.Username, user.LangCode, joinedAt, deactivatedAt, rejoinedAt))
}

func (p *Processor) rejoined(ctx context.Context, user *tg.User) {

	postgres.RejoinUser(ctx, user)
	joinedAt := format(postgres.JoinedAt(ctx, user))
	deactivatedAt := format(postgres.DeactivatedAt(ctx, user))
	rejoinedAt := format(postgres.RejoinedAt(ctx, user))

	if len(user.Username) > 0 {
		user.Username = "@" + user.Username
	}

	info.SendInfoMessage(ctx, fmt.Sprintf(
		`is_active: ✅

id: %d
first_name: %s
last_name: %s
username: %s
language_code: %s

joined_at:           %s
deactivated_at: %s
rejoined_at:       %s`, user.ID, user.FirstName, user.LastName, user.Username, user.LangCode, joinedAt, deactivatedAt, rejoinedAt))
}
