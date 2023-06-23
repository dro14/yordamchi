package redis

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/types"
)

func LoadContext(ctx context.Context, prompt string) []types.Message {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	messages := make([]types.Message, 0, 3)
	message := types.Message{
		Content: prompt,
		Role:    "user",
	}

	result, err := Client.Get(ctx, "context:"+id).Result()
	if err != nil {
		if err.Error() != e.KeyNotFound {
			log.Printf("can't get \"context:%s\": %v", id, err)
		}
		return append(messages, message)
	}

	err = json.Unmarshal([]byte(result), &messages)
	if err != nil {
		log.Printf("can't decode context: %v", err)
		return append(messages, message)
	}

	return append(messages, message)
}

func StoreContext(ctx context.Context, prompt, completion string) {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	messages := []types.Message{
		{
			Content: prompt,
			Role:    "user",
		},
		{
			Content: completion,
			Role:    "assistant",
		},
	}

	data, err := json.Marshal(messages)
	if err != nil {
		log.Printf("can't encode context: %v", err)
		return
	}

	err = Client.Set(ctx, "context:"+id, string(data), 72*time.Hour).Err()
	if err != nil {
		log.Printf("can't store context: %v", err)
	}
}

func DeleteContext(ctx context.Context) {

	id := strconv.Itoa(int(ctx.Value("user_id").(int64)))

	err := Client.Del(ctx, "context:"+id).Err()
	if err != nil {
		log.Printf("can't delete context: %v", err)
	}
}
