package cache

import (
	"context"

	"github.com/dro14/yordamchi/lib/types"
	"github.com/gotd/td/tg"
)

type Cache interface {
	Status(context.Context) (types.UserStatus, error)
	Balance(context.Context) (int, error)
	Decrement(context.Context) error
	LoadContext(context.Context, string) []types.Message
	StoreContext(context.Context, string, string)
	DeleteContext(context.Context)
	IncrementActivity(context.Context, *tg.Message, *tg.User, bool) int
	DecrementActivity(context.Context)
	LoadActivity(context.Context) []*types.Activity
}
