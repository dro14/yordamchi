package redis

import (
	"log"
	"os"

	"github.com/dro14/yordamchi/clients/other"
	"github.com/dro14/yordamchi/clients/service"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client  *redis.Client
	service *service.Service
	apis    *other.APIs
}

func New() *Redis {
	url, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		log.Fatal("redis url is not specified")
	}

	password, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		log.Fatal("redis password is not specified")
	}

	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     url,
			Password: password,
		}),
		service: service.New(),
		apis:    other.New(),
	}
}
