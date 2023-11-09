package redis

import (
	"context"
)

func Lang(ctx context.Context, languageCode string) (context.Context, bool) {
	switch languageCode {
	case "uz", "":
		ctx = context.WithValue(ctx, "language_code", "uz")
	case "ru":
		ctx = context.WithValue(ctx, "language_code", "ru")
	default:
		ctx = context.WithValue(ctx, "language_code", "en")
	}
	lang, err := Client.Get(ctx, "lang:"+id(ctx)).Result()
	if err != nil {
		return ctx, false
	}
	ctx = context.WithValue(ctx, "language_code", lang)
	return ctx, true
}

func SetLang(ctx context.Context) {
	lang := ctx.Value("language_code").(string)
	Client.Set(ctx, "lang:"+id(ctx), lang, 0)
}
