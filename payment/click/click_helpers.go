package click

import "context"

func id(ctx context.Context) int64 {
	return ctx.Value("user_id").(int64)
}
