package redis

import (
	"context"
	"fmt"
)

func id(ctx context.Context) string {
	return fmt.Sprintf("%d", ctx.Value("user_id").(int64))
}

func userStatus(ctx context.Context) UserStatus {
	return ctx.Value("user_status").(UserStatus)
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

func model(ctx context.Context) string {
	return ctx.Value("model").(string)
}

func translate(ctx context.Context) bool {
	return ctx.Value("translate").(bool)
}
