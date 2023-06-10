package database

import (
	"context"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/gotd/td/tg"
)

type Database interface {
	JoinUser(context.Context, *tg.User, int64)
	SaveMessage(context.Context, *types.Stats, *tg.User)
	DeactivateUser(context.Context, *tg.User)
	RejoinUser(context.Context, *tg.User)
	IsActive(context.Context, *tg.User) bool
	JoinedAt(context.Context, *tg.User) string
	DeactivatedAt(context.Context, *tg.User) string
	RejoinedAt(context.Context, *tg.User) string
}
