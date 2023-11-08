package redis

import (
	"context"
	"fmt"
)

func Lang(ctx context.Context, languageCode string) (context.Context, bool) {
	if languageCode == "" {
		ctx = context.WithValue(ctx, "language_code", "uz")
	} else if languageCode != "uz" && languageCode != "ru" {
		ctx = context.WithValue(ctx, "language_code", "en")
	}

	key := fmt.Sprintf("lang:%d", ctx.Value("user_id").(int64))
	lang, err := Client.Get(ctx, key).Result()
	if err != nil {
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
