package redis

import (
	"context"

	"github.com/dro14/yordamchi/lib/models"
)

func GPT3(ctx context.Context) {
	Client.Del(ctx, "model:"+id(ctx))
}

func GPT4(ctx context.Context) {
	Client.Set(ctx, "model:"+id(ctx), models.GPT4, 0)
}

func Model(ctx context.Context) string {
	_, err := Client.Get(ctx, "model:"+id(ctx)).Result()
	if err != nil {
		return models.GPT3
	}
	return models.GPT4
}

func GPT4Tokens(ctx context.Context) int {
	tokens, _ := Client.Get(ctx, "gpt-4:"+id(ctx)).Int()
	return tokens
}
