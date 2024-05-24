package openai

import (
	"fmt"
	"log"
	"os"

	"github.com/dro14/yordamchi/clients/other"
	"github.com/dro14/yordamchi/clients/service"
	"github.com/dro14/yordamchi/storage/redis"
	"github.com/pkoukk/tiktoken-go"
)

type OpenAI struct {
	keys    []string
	index   int
	tkm     *tiktoken.Tiktoken
	redis   *redis.Redis
	service *service.Service
	apis    *other.APIs
}

func New() *OpenAI {
	var keys []string
	for i := 0; ; i++ {
		key := fmt.Sprintf("OPENAI_API_KEY_%d", i)
		token, ok := os.LookupEnv(key)
		if !ok {
			break
		}
		keys = append(keys, "Bearer "+token)
	}

	if len(keys) == 0 {
		log.Fatal("openai token is not specified")
	}

	tkm, err := tiktoken.GetEncoding(tiktoken.MODEL_CL100K_BASE)
	if err != nil {
		log.Fatal("can't get encoding:", err)
	}

	return &OpenAI{
		keys:    keys,
		index:   0,
		tkm:     tkm,
		redis:   redis.New(),
		service: service.New(),
		apis:    other.New(),
	}
}
