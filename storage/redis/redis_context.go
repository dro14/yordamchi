package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dro14/yordamchi/clients/openai/models"
	"github.com/dro14/yordamchi/clients/openai/types"
	"github.com/dro14/yordamchi/utils"
)

func (r *Redis) Context(ctx context.Context, prompt string) (context.Context, []types.Message) {
	searchTemplate := map[string]string{
		"uz": "\n\nQUYIDA MAVZUGA OID MA'LUMOTLAR KELTIRILGAN. KERAK BO'LSA ULARDAN FOYDALAN.\n\n%s",
		"ru": "\n\nНИЖЕ ПРИВЕДЕНЫ СООТВЕТСТВУЮЩИЕ ТЕМЕ ФРАГМЕНТЫ ИНФОРМАЦИИ. ИСПОЛЬЗУЙ ИХ, ЕСЛИ ОНИ БУДУТ ПОЛЕЗНЫ.\n\n%s",
		"en": "\n\nTHE FOLLOWING ARE THE RELEVANT PIECES OF INFORMATION. USE THEM IF HELPFUL.\n\n%s",
	}

	system := r.System(ctx)
	messages := r.messages(ctx)
	if userStatus(ctx) != StatusFree && !strings.Contains(prompt, utils.Delim) {
		var query string
		if len(messages) == 2 && !strings.Contains(messages[0].Content.(string), utils.Delim) {
			query = messages[0].Content.(string) + " " + prompt
		} else {
			query = prompt
		}

		if model(ctx) == models.GPT3 && lang(ctx) == "uz" {
			system += fmt.Sprintf(searchTemplate["en"], r.service.Search(ctx, query, "en"))
		} else {
			system += fmt.Sprintf(searchTemplate[lang(ctx)], r.service.Search(ctx, query, lang(ctx)))
		}
	}

	system, _ = strings.CutPrefix(system, "USER: ")
	messages = append([]types.Message{{Role: "system", Content: system}}, messages...)
	messages = append(messages, types.Message{Role: "user", Content: prompt})

	for i := range messages {
		URL, text, found := strings.Cut(messages[i].Content.(string), utils.Delim)
		if found {
			var content []types.Content
			if text != "" {
				content = append(content, types.Content{Type: "text", Text: text})
			}
			content = append(content, types.Content{Type: "image_url", ImageURL: &types.ImageURL{URL: URL}})
			messages[i].Content = content
			ctx = context.WithValue(ctx, "model", models.GPT4V)
		}
	}
	return ctx, messages
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

	var expiration time.Duration
	if strings.Contains(prompt, utils.Delim) {
		expiration = 1 * time.Hour
	} else {
		expiration = 72 * time.Hour
	}
	r.client.Set(ctx, "context:"+id(ctx), jsonData, expiration)
}

func (r *Redis) DeleteContext(ctx context.Context) {
	r.client.Del(ctx, "context:"+id(ctx))
	r.client.Del(ctx, "system:"+id(ctx))
	err := r.service.Delete(ctx)
	if err != nil {
		log.Printf("user %s: can't delete file", id(ctx))
	}
}

func (r *Redis) System(ctx context.Context) string {
	promptTemplate := map[string]string{
		"uz": "SEN %s MODEL ARXITEKTURASIGA ASOSLANGAN, TELEGRAMDAGI YORDAMCHI NOMLI XUSHMUOMALA CHATBOTSAN.",
		"ru": "ТЫ ЯВЛЯЕШЬСЯ ДРУЖЕЛЮБНЫМ ЧАТБОТОМ В ТЕЛЕГРАМЕ ПОД НАЗВАНИЕМ YORDAMCHI, ОСНОВАННЫЙ НА АРХИТЕКТУРЕ МОДЕЛИ %s.",
		"en": "YOU ARE A FRIENDLY CHATBOT IN TELEGRAM CALLED YORDAMCHI, BASED ON %s MODEL ARCHITECTURE.",
	}

	system, err := r.client.Get(ctx, "system:"+id(ctx)).Result()
	if err != nil || userStatus(ctx) == StatusFree {
		if model(ctx) == models.GPT3 {
			if lang(ctx) == "uz" {
				return fmt.Sprintf(promptTemplate["en"], "GPT-3.5")
			} else {
				return fmt.Sprintf(promptTemplate[lang(ctx)], "GPT-3.5")
			}
		} else {
			return fmt.Sprintf(promptTemplate[lang(ctx)], "GPT-4")
		}
	}
	return "USER: " + system
}

func (r *Redis) SetSystem(ctx context.Context, system string) {
	r.client.Set(ctx, "system:"+id(ctx), system, 0)
}

func (r *Redis) messages(ctx context.Context) []types.Message {
	jsonData, err := r.client.Get(ctx, "context:"+id(ctx)).Bytes()
	if err != nil {
		return []types.Message{}
	}

	var messages []types.Message
	err = json.Unmarshal(jsonData, &messages)
	if err != nil {
		log.Printf("can't decode %q: %s", "context:"+id(ctx), err)
	}
	return messages
}
