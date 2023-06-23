package telegram

import (
	"context"
	"log"

	"github.com/dro14/yordamchi/client/telegram"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/processor/openai"
	"github.com/dro14/yordamchi/redis"
	"github.com/gotd/td/tg"
)

type Processor struct {
	Client    *telegram.Client
	Processor *openai.Processor
}

func New(bot *tg.Client) *Processor {
	return &Processor{
		Client:    telegram.New(bot),
		Processor: openai.New(),
	}
}

func (p *Processor) ProcessMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) error {

	ctx, message, user := messageUpdate(ctx, entities, update)
	if message == nil && user == nil {
		return nil
	}
	defer blockedUsers.Delete(user.ID)

	done := p.doCommand(ctx, message, user)
	if done {
		return nil
	}

	userStatus, err := redis.Status(ctx)

	switch userStatus {
	case types.UnknownStatus:
		log.Printf("unknown user status: %v", err)
	case types.BlockedStatus:
		p.blocked(ctx)
	case types.PremiumStatus:
		p.Stream(ctx, message, user, true)
	case types.FreeStatus:
		p.Stream(ctx, message, user, false)
	case types.ExhaustedStatus:
		p.exhausted(ctx)
	}

	return nil
}

func (p *Processor) ProcessCallbackQuery(ctx context.Context, entities tg.Entities, update *tg.UpdateBotCallbackQuery) error {

	ctx, data := callbackUpdate(ctx, entities, update)

	switch data {
	case "new_chat":
		p.newChatCallback(ctx)
	case "examples":
		p.examplesCallback(ctx, update.MsgID)
	case "help":
		p.helpCallback(ctx, update.MsgID)
	case "premium":
		p.premiumCallback(ctx, update.MsgID)
	default:
		p.confirmCallback(ctx, update.MsgID, data)
	}

	return nil
}

func (p *Processor) ProcessBotStopped(ctx context.Context, entities tg.Entities, update *tg.UpdateBotStopped) error {

	ctx, user := botStoppedUpdate(ctx, entities, update)

	if update.Stopped {
		p.deactivated(ctx, user)
	} else {
		p.rejoined(ctx, user)
	}

	return nil
}
