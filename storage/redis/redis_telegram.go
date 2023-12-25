package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dro14/yordamchi/utils"
	"github.com/go-redis/redis/v8"
)

type UserStatus int

const (
	StatusUnknown UserStatus = iota
	StatusExhausted
	StatusFree
	StatusUnlimited
	StatusPremium
)

func (r *Redis) UserStatus(ctx context.Context) UserStatus {
	_, err := r.client.Get(ctx, "premium:"+id(ctx)).Result()
	if err == nil {
		return StatusPremium
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("user %s: can't check whether status is premium: %s", id(ctx), err)
		return StatusUnknown
	}

	_, err = r.client.Get(ctx, "unlimited:"+id(ctx)).Result()
	if err == nil {
		return StatusUnlimited
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("user %s: can't check whether status is unlimited: %s", id(ctx), err)
		return StatusUnknown
	}

	requests, err := r.client.Get(ctx, "free:"+id(ctx)).Int()
	if err == nil {
		if requests > 0 && requests <= utils.NumOfFreeReqs {
			return StatusFree
		} else if requests != 0 {
			log.Printf("user %s: invalid number of requests: %d", id(ctx), requests)
			return StatusUnknown
		}
	} else if errors.Is(err, redis.Nil) {
		r.client.Set(ctx, "free:"+id(ctx), utils.NumOfFreeReqs, untilMidnight())
		return StatusFree
	} else {
		log.Printf("user %s: can't check whether status is free: %s", id(ctx), err)
		return StatusUnknown
	}

	return StatusExhausted
}

func (r *Redis) Expiration(ctx context.Context) string {
	expirationDate, err := r.client.Get(ctx, "premium:"+id(ctx)).Result()
	if err == nil {
		return expirationDate
	}

	expirationDate, err = r.client.Get(ctx, "unlimited:"+id(ctx)).Result()
	if err == nil {
		return expirationDate
	}

	return midnight()
}

func (r *Redis) Requests(ctx context.Context) string {
	requests, err := r.client.Get(ctx, "free:"+id(ctx)).Int()
	if err != nil {
		log.Printf("can't get %q: %v", "free:"+id(ctx), err)
		return ""
	}
	return fmt.Sprintf("%d/%d", requests, utils.NumOfFreeReqs)
}

func (r *Redis) DecrementRequests(ctx context.Context) {
	if ctx.Value("user_status") == StatusFree {
		requests, err := r.client.Get(ctx, "free:"+id(ctx)).Int()
		if err != nil {
			log.Printf("can't get %q: %s", "free:"+id(ctx), err)
			return
		}
		if requests > 0 && requests <= utils.NumOfFreeReqs {
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
	langCode, err := r.client.Get(ctx, "lang:"+id(ctx)).Result()
	if err != nil {
		return ctx, false
	}
	ctx = context.WithValue(ctx, "language_code", langCode)
	r.SetLang(ctx)
	return ctx, true
}

func (r *Redis) SetLang(ctx context.Context) {
	expiration := time.Now().AddDate(0, 1, 0)
	r.client.Set(ctx, "lang:"+id(ctx), lang(ctx), time.Until(expiration))
}

func (r *Redis) PollQuestion(ctx context.Context) string {
	return r.client.Get(ctx, "poll_question").String()
}
