package redis

import (
	"context"
	"fmt"
	"log"
)

func Lang(ctx context.Context, languageCode string) (context.Context, error) {

	if languageCode == "" {
		languageCode = "uz"
	} else if languageCode != "uz" && languageCode != "ru" {
		languageCode = "en"
	}
	ctx = context.WithValue(ctx, "language_code", languageCode)

	key := fmt.Sprintf("tl:%d", ctx.Value("user_id").(int64))
	lang, err := Client.Get(ctx, key).Result()
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, "target_lang", lang)
	return ctx, nil
}

func SetLang(ctx context.Context, lang string) {

	key := fmt.Sprintf("tl:%d", ctx.Value("user_id").(int64))

	err := Client.Set(ctx, key, lang, 0).Err()
	if err != nil {
		log.Printf("can't set %q: %v", key, err)
	}
}
