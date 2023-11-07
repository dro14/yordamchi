package redis

import (
	"context"
	"fmt"
	"github.com/dro14/yordamchi/lib/models"
)

func GPT3(ctx context.Context) error {
	key := fmt.Sprintf("model:%d", ctx.Value("user_id").(int64))
	return Client.Del(ctx, key).Err()
}

func GPT4(ctx context.Context) error {
	key := fmt.Sprintf("model:%d", ctx.Value("user_id").(int64))
	return Client.Set(ctx, key, models.GPT4, 0).Err()
}

func Model(ctx context.Context) string {
	key := fmt.Sprintf("model:%d", ctx.Value("user_id").(int64))
	_, err := Client.Get(ctx, key).Result()
	if err != nil {
		return models.GPT3
	}
	return models.GPT4
}

func GPT4Tokens(ctx context.Context) int {
	key := fmt.Sprintf("gpt-4:%d", ctx.Value("user_id").(int64))
	tokens, err := Client.Get(ctx, key).Int()
	if err != nil {
		return 0
	}
	return tokens
}
