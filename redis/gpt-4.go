package redis

import (
	"context"
	"fmt"
)

func GPT3(ctx context.Context) error {
	key := fmt.Sprintf("model:%d", ctx.Value("user_id").(int64))
	return Client.Del(ctx, key).Err()
}

func GPT4(ctx context.Context) error {
	key := fmt.Sprintf("model:%d", ctx.Value("user_id").(int64))
	return Client.Set(ctx, key, "gpt-4", 0).Err()
}

func Model(ctx context.Context) string {
	key := fmt.Sprintf("model:%d", ctx.Value("user_id").(int64))
	model, err := Client.Get(ctx, key).Result()
	if err != nil {
		return "gpt-3.5-turbo"
	}
	return model
}

func GPT4Tokens(ctx context.Context) int {
	key := fmt.Sprintf("gpt-4:%d", ctx.Value("user_id").(int64))
	tokens, err := Client.Get(ctx, key).Int()
	if err != nil {
		return 0
	}
	return tokens
}
