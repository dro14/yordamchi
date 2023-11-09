package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/dro14/yordamchi/lib/types"
)

func LoadContext(ctx context.Context, prompt string) []types.Message {
	var messages []types.Message
	message := types.Message{Content: prompt, Role: "user"}

	jsonData, err := Client.Get(ctx, "context:"+id(ctx)).Result()
	if err != nil {
		return append(messages, message)
	}

	err = json.Unmarshal([]byte(jsonData), &messages)
	if err != nil {
		log.Printf("can't decode %q: %v", "context:"+id(ctx), err)
		return append(messages, message)
	}
	return append(messages, message)
}

func StoreContext(ctx context.Context, prompt, completion string) {
	messages := []types.Message{
		{Content: prompt, Role: "user"},
		{Content: completion, Role: "assistant"},
	}
	jsonData, err := json.Marshal(messages)
	if err != nil {
		log.Printf("can't encode %q: %v", "context:"+id(ctx), err)
		return
	}
	Client.Set(ctx, "context:"+id(ctx), string(jsonData), 72*time.Hour)
}

func DeleteContext(ctx context.Context) {
	Client.Del(ctx, "context:"+id(ctx))
}
