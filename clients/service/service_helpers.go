package service

import "context"

func id(ctx context.Context) int64 {
	return ctx.Value("user_id").(int64)
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}
