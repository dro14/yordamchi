package redis

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/clients/openai/types"
	"github.com/dro14/yordamchi/utils"
)

var template = map[string]string{
	"uz": "SEN TELEGRAMDAGI YORDAMCHI NOMLI XUSHMUOMALA CHATBOTSAN.",
	"ru": "ТЫ ЯВЛЯЕШЬСЯ ДРУЖЕЛЮБНЫМ ЧАТБОТОМ В ТЕЛЕГРАМЕ ПОД НАЗВАНИЕМ YORDAMCHI.",
	"en": "YOU ARE A FRIENDLY CHATBOT IN TELEGRAM CALLED YORDAMCHI.",
}

func (r *Redis) Context(ctx context.Context, prompt *string) []types.Message {
	system := r.System(ctx)
	system, _ = strings.CutPrefix(system, "USER: ")
	messages := r.Messages(ctx)

	if translate(ctx) {
		*prompt = r.apis.Translate("uz", "en", *prompt)
	}

	messages = append([]types.Message{{Role: "system", Content: system}}, messages...)
	messages = append(messages, types.Message{Role: "user", Content: *prompt})

	for i := range messages {
		URL, text, found := strings.Cut(messages[i].Content.(string), utils.Delim)
		if found {
			var content []types.Content
			if text != "" {
				content = append(content, types.Content{Type: "text", Text: text})
			}
			content = append(content, types.Content{
				Type:     "image_url",
				ImageURL: &types.ImageURL{URL: URL, Detail: "low"},
			})
			messages[i].Content = content
		}
	}
	return messages
}

func (r *Redis) SetContext(ctx context.Context, prompt, completion string) {
	messages := []types.Message{
		{Role: "user", Content: prompt},
		{Role: "assistant", Content: completion},
	}

	jsonData, err := json.Marshal(messages)
	if err != nil {
		log.Printf("can't encode %q: %s", "context:"+id(ctx), err)
		return
	}

	expiration := time.Now().UTC()
	if strings.Contains(prompt, utils.Delim) {
		expiration = expiration.Add(1 * time.Hour)
	} else {
		expiration = expiration.Add(72 * time.Hour)
	}
	client.Set(ctx, "context:"+id(ctx), jsonData, time.Until(expiration))
}

func (r *Redis) DeleteContext(ctx context.Context) {
	client.Del(ctx, "context:"+id(ctx))
	client.Del(ctx, "system:"+id(ctx))
	r.service.Delete(ctx)
}

func (r *Redis) System(ctx context.Context) string {
	system, err := client.Get(ctx, "system:"+id(ctx)).Result()
	if err != nil || userStatus(ctx) == StatusFree {
		if translate(ctx) {
			return template["en"]
		} else {
			return template[lang(ctx)]
		}
	}
	return "USER: " + system
}

func (r *Redis) SetSystem(ctx context.Context, system string) {
	client.Set(ctx, "system:"+id(ctx), system, 0)
}

func (r *Redis) Messages(ctx context.Context) []types.Message {
	jsonData, err := client.Get(ctx, "context:"+id(ctx)).Bytes()
	if err != nil {
		return nil
	}

	var messages []types.Message
	err = json.Unmarshal(jsonData, &messages)
	if err != nil {
		log.Printf("can't decode %q: %s", "context:"+id(ctx), err)
		return nil
	}

	if model(ctx) == models.GPT3 && strings.Contains(messages[0].Content.(string), utils.Delim) {
		return nil
	} else {
		return messages
	}
}
