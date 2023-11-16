package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/dro14/yordamchi/utils"
	"log"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/go-redis/redis/v8"
)

func (r *Redis) UserStatus(ctx context.Context) UserStatus {
	_, err := r.client.Get(ctx, "model:"+id(ctx)).Result()
	if err == nil {
		return GPT4Status
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("user %s: can't check whether status is gpt-4: %s", id(ctx), err)
		return UnknownStatus
	}

	_, err = r.client.Get(ctx, "premium:"+id(ctx)).Result()
	if err == nil {
		return PremiumStatus
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("user %s: can't check whether status is premium: %s", id(ctx), err)
		return UnknownStatus
	}

	requests, err := r.client.Get(ctx, "free:"+id(ctx)).Int()
	if err == nil {
		if requests > 0 && requests <= utils.NumOfFreeRequests {
			return FreeStatus
		} else if requests != 0 {
			log.Printf("user %s: invalid number of requests: %d", id(ctx), requests)
			return UnknownStatus
		}
	} else if errors.Is(err, redis.Nil) {
		r.client.Set(ctx, "free:"+id(ctx), utils.NumOfFreeRequests, untilMidnight())
		return FreeStatus
	} else {
		log.Printf("user %s: can't check whether status is free: %s", id(ctx), err)
		return UnknownStatus
	}

	return ExhaustedStatus
}

func (r *Redis) Expiration(ctx context.Context) string {
	expirationDate, err := r.client.Get(ctx, "premium:"+id(ctx)).Result()
	if err != nil {
		return midnight()
	}
	return expirationDate
}

func (r *Redis) Requests(ctx context.Context) string {
	requests, err := r.client.Get(ctx, "free:"+id(ctx)).Int()
	if err != nil {
		log.Printf("can't get %q: %v", "free:"+id(ctx), err)
		return ""
	}
	return fmt.Sprintf("%d/%d", requests, utils.NumOfFreeRequests)
}

func (r *Redis) Decrement(ctx context.Context, used int) {
	switch ctx.Value("model") {
	case models.GPT4, models.GPT4V:
		available, err := r.client.Get(ctx, "gpt-4:"+id(ctx)).Int()
		if err != nil {
			log.Printf("can't get %q: %s", "gpt-4:"+id(ctx), err)
			return
		}
		if available <= used {
			r.client.Del(ctx, "gpt-4:"+id(ctx))
		} else {
			r.client.Set(ctx, "gpt-4:"+id(ctx), available-used, 0)
		}
	default:
		_, err := r.client.Get(ctx, "premium:"+id(ctx)).Result()
		if err == nil {
			return
		}
		requests, err := r.client.Get(ctx, "free:"+id(ctx)).Int()
		if err != nil {
			log.Printf("can't get %q: %s", "free:"+id(ctx), err)
			return
		}
		if requests > 0 && requests <= utils.NumOfFreeRequests {
			r.client.Set(ctx, "free:"+id(ctx), requests-1, untilMidnight())
		} else {
			log.Printf("user %s: invalid number of requests: %d", id(ctx), requests)
		}
	}
}

func (r *Redis) Lang(ctx context.Context, languageCode string) (context.Context, bool) {
	switch languageCode {
	case "uz", "":
		ctx = context.WithValue(ctx, "language_code", "uz")
	case "ru":
		ctx = context.WithValue(ctx, "language_code", "ru")
	default:
		ctx = context.WithValue(ctx, "language_code", "en")
	}
	lang, err := r.client.Get(ctx, "lang:"+id(ctx)).Result()
	if err != nil {
		return ctx, false
	}
	ctx = context.WithValue(ctx, "language_code", lang)
	return ctx, true
}

func (r *Redis) SetLang(ctx context.Context) {
	lang := ctx.Value("language_code").(string)
	r.client.Set(ctx, "lang:"+id(ctx), lang, 0)
}

func (r *Redis) GPT3(ctx context.Context) {
	r.client.Del(ctx, "model:"+id(ctx))
}

func (r *Redis) GPT4(ctx context.Context) {
	r.client.Set(ctx, "model:"+id(ctx), models.GPT4, 0)
}

func (r *Redis) Model(ctx context.Context) string {
	_, err := r.client.Get(ctx, "model:"+id(ctx)).Result()
	if err != nil {
		return models.GPT3
	}
	return models.GPT4
}

func (r *Redis) GPT4Tokens(ctx context.Context) int {
	tokens, _ := r.client.Get(ctx, "gpt-4:"+id(ctx)).Int()
	return tokens
}
