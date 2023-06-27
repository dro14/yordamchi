package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dro14/yordamchi/lib/types"
)

func LoadContext(ctx context.Context, prompt string) []types.Message {

	key := fmt.Sprintf("context:%d", ctx.Value("user_id").(int64))

	messages := make([]types.Message, 0, 3)
	message := types.Message{
		Content: prompt,
		Role:    "user",
	}

	result, err := Client.Get(ctx, key).Result()
	if err != nil {
		if err.Error() != KeyNotFound {
			log.Printf("can't get %q: %v", key, err)
		}
		return append(messages, message)
	}

	err = json.Unmarshal([]byte(result), &messages)
	if err != nil {
		log.Printf("can't decode %q: %v", key, err)
		return append(messages, message)
	}

	return append(messages, message)
}

func StoreContext(ctx context.Context, prompt, completion string) {

	key := fmt.Sprintf("context:%d", ctx.Value("user_id").(int64))

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
		log.Printf("can't encode %q: %v", key, err)
		return
	}

	err = Client.Set(ctx, key, string(data), 72*time.Hour).Err()
	if err != nil {
		log.Printf("can't set %q: %v", key, err)
	}
}

func DeleteContext(ctx context.Context) {

	key := fmt.Sprintf("context:%d", ctx.Value("user_id").(int64))

	err := Client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("can't delete %q: %v", key, err)
	}
}
