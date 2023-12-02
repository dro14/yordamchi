package redis

import (
	"context"
	"log"
)

func (r *Redis) Images(ctx context.Context) int {
	images, _ := r.client.Get(ctx, "images:"+id(ctx)).Int()
	return images
}

func (r *Redis) DecrementImages(ctx context.Context) {
	images, err := r.client.Get(ctx, "images:"+id(ctx)).Int()
	if err != nil {
		log.Printf("can't get %q: %s", "images:"+id(ctx), err)
		return
	}

	if images > 1 {
		r.client.Set(ctx, "images:"+id(ctx), images-1, 0)
	} else if images == 1 {
		r.client.Del(ctx, "images:"+id(ctx))
	} else {
		log.Printf("user %s: invalid number of images: %d", id(ctx), images)
	}
}

func (r *Redis) Generate(ctx context.Context) string {
	generate, err := r.client.Get(ctx, "generate:"+id(ctx)).Result()
	if err != nil {
		log.Printf("can't get %q: %s", "generate:"+id(ctx), err)
		return ""
	}
	r.client.Del(ctx, "generate:"+id(ctx))
	return generate
}

func (r *Redis) SetGenerate(ctx context.Context, generate string) {
	r.client.Set(ctx, "generate:"+id(ctx), generate, 0)
}
