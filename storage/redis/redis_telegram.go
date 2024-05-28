package redis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
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
	_, err := client.Get(ctx, "premium:"+id(ctx)).Result()
	if err == nil {
		return StatusPremium
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("user %s: can't check whether status is premium: %s", id(ctx), err)
		return StatusUnknown
	}

	_, err = client.Get(ctx, "unlimited:"+id(ctx)).Result()
	if err == nil {
		return StatusUnlimited
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("user %s: can't check whether status is unlimited: %s", id(ctx), err)
		return StatusUnknown
	}

	requests, err := client.Get(ctx, "free:"+id(ctx)).Int()
	if err == nil {
		if requests > 0 {
			return StatusFree
		} else if requests != 0 {
			log.Printf("user %s: invalid number of requests: %d", id(ctx), requests)
			return StatusUnknown
		}
	} else if errors.Is(err, redis.Nil) {
		client.Set(ctx, "free:"+id(ctx), utils.NumOfFreeReqs, 0)
		return StatusFree
	} else {
		log.Printf("user %s: can't check whether status is free: %s", id(ctx), err)
		return StatusUnknown
	}

	return StatusExhausted
}

func (r *Redis) Expiration(ctx context.Context) string {
	expiration, err := client.Get(ctx, "premium:"+id(ctx)).Result()
	if err == nil {
		return strings.Split(expiration, "|")[0]
	}

	expiration, err = client.Get(ctx, "unlimited:"+id(ctx)).Result()
	if err == nil {
		return expiration
	}

	log.Printf("can't get expiration for %s: %v", id(ctx), err)
	return ""
}

func (r *Redis) Requests(ctx context.Context) string {
	requests, err := client.Get(ctx, "premium:"+id(ctx)).Result()
	if err == nil {
		return strings.Split(requests, "|")[1]
	}

	requests, err = client.Get(ctx, "free:"+id(ctx)).Result()
	if err == nil {
		return requests
	}

	log.Printf("can't get requests for %s: %v", id(ctx), err)
	return ""
}

func (r *Redis) Premium(ctx context.Context) (string, string) {
	value, err := client.Get(ctx, "premium:"+id(ctx)).Result()
	if err != nil {
		log.Printf("can't get %q: %s", "premium:"+id(ctx), err)
		return "", ""
	}

	values := strings.Split(value, "|")
	return values[0], values[1]
}

func (r *Redis) DecrementRequests(ctx context.Context) {
	if ctx.Value("user_status") == StatusPremium {
		value, err := client.Get(ctx, "premium:"+id(ctx)).Result()
		if err != nil {
			log.Printf("can't get %q: %s", "premium:"+id(ctx), err)
			return
		}

		values := strings.Split(value, "|")
		requests, _ := strconv.Atoi(values[1])
		value = fmt.Sprintf("%s|%d", values[0], requests-1)
		expiration, _ := time.Parse("02.01.2006 15:04:05", values[0])

		if requests > 1 {
			client.Set(ctx, "premium:"+id(ctx), value, time.Until(expiration))
		} else if requests == 1 {
			client.Del(ctx, "premium:"+id(ctx))
		} else {
			log.Printf("user %s: invalid number of requests: %d", id(ctx), requests)
		}
	} else if ctx.Value("user_status") == StatusFree {
		requests, err := client.Get(ctx, "free:"+id(ctx)).Int()
		if err != nil {
			log.Printf("can't get %q: %s", "free:"+id(ctx), err)
			return
		}
		if requests > 0 {
			client.Set(ctx, "free:"+id(ctx), requests-1, 0)
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
	langCode, err := client.Get(ctx, "lang:"+id(ctx)).Result()
	if err != nil {
		return ctx, false
	}
	ctx = context.WithValue(ctx, "language_code", langCode)
	r.SetLang(ctx)
	return ctx, true
}

func (r *Redis) SetLang(ctx context.Context) {
	expiration := time.Now().AddDate(0, 1, 0)
	client.Set(ctx, "lang:"+id(ctx), lang(ctx), time.Until(expiration))
}

func (r *Redis) PollQuestion(ctx context.Context) string {
	return client.Get(ctx, "poll_question").Val()
}
