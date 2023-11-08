package redis

import (
	"context"
	"fmt"

	"github.com/dro14/yordamchi/lib/e"
)

func Lang(ctx context.Context) (context.Context, bool) {
	key := fmt.Sprintf("lang:%d", ctx.Value("user_id").(int64))
	lang, err := Client.Get(ctx, key).Result()
	if err.Error() == e.KeyNotFound {
		return ctx, true
	}
	ctx = context.WithValue(ctx, "language_code", lang)
	return ctx, false
}

func SetLang(ctx context.Context) {
	key := fmt.Sprintf("lang:%d", ctx.Value("user_id").(int64))
	lang := ctx.Value("language_code").(string)
	Client.Set(ctx, key, lang, 0)
}
