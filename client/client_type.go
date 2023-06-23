package client

import (
	"context"
	"sync/atomic"

	"github.com/dro14/yordamchi/lib/types"
	"github.com/gotd/td/tg"
)

type TelegramClient interface {
	SendMessage(context.Context, string, int, *tg.ReplyInlineMarkup) (int, error)
	EditMessage(context.Context, string, int, *tg.ReplyInlineMarkup) error
	SetTyping(context.Context, *atomic.Bool)
}

type OpenAIClient interface {
	Completion(context.Context, []types.Message, int, chan<- string) (*types.Response, error)
}
