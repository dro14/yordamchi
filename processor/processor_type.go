package processor

import (
	"context"

	"github.com/dro14/yordamchi/lib/types"
	"github.com/gotd/td/tg"
)

type TelegramProcessor interface {
	ProcessMessage(ctx context.Context, entities tg.Entities, update *tg.UpdateNewMessage) error
	ProcessCallbackQuery(ctx context.Context, entities tg.Entities, update *tg.UpdateBotCallbackQuery) error
	ProcessBotStopped(ctx context.Context, entities tg.Entities, update *tg.UpdateBotStopped) error
}

type OpenAIProcessor interface {
	Process(ctx context.Context, messages []types.Message, stats *types.Stats, channel chan<- string)
}
