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

var promptTemplate = map[string]string{
	"uz": "SEN %s MODEL ARXITEKTURASIGA ASOSLANGAN, TELEGRAMDAGI YORDAMCHI NOMLI XUSHMUOMALA CHATBOTSAN.",
	"ru": "ТЫ ЯВЛЯЕШЬСЯ ДРУЖЕЛЮБНЫМ ЧАТБОТОМ В ТЕЛЕГРАМЕ ПОД НАЗВАНИЕМ YORDAMCHI, ОСНОВАННЫЙ НА АРХИТЕКТУРЕ МОДЕЛИ %s.",
	"en": "YOU ARE A FRIENDLY CHATBOT IN TELEGRAM CALLED YORDAMCHI, BASED ON %s MODEL ARCHITECTURE.",
}

var searchTemplate = map[string]string{
	"uz": " QUYIDA MAVZUGA OID MA'LUMOTLAR KELTIRILGAN. KERAK BO'LSA ULARDAN FOYDALAN.\n\n%s",
	"ru": " НИЖЕ ПРИВЕДЕНЫ СООТВЕТСТВУЮЩИЕ ТЕМЕ ФРАГМЕНТЫ ИНФОРМАЦИИ. ИСПОЛЬЗУЙ ИХ, ЕСЛИ ОНИ БУДУТ ПОЛЕЗНЫ.\n\n%s",
	"en": " THE FOLLOWING ARE THE RELEVANT PIECES OF INFORMATION. USE THEM IF HELPFUL.\n\n%s",
}

func (r *Redis) ConversationHistory(ctx context.Context, prompt string) (output context.Context, messages []types.Message) {
	jsonData, err := r.client.Get(ctx, "context:"+id(ctx)).Bytes()
	if err == nil {
		err = json.Unmarshal(jsonData, &messages)
		if err != nil {
			log.Printf("can't decode %q: %s", "context:"+id(ctx), err)
		}
	}

	var sysPrompt string
	if ctx.Value("model") == models.GPT3 {
		if lang(ctx) == "uz" {
			sysPrompt = fmt.Sprintf(promptTemplate["en"], "GPT-3.5")
		} else {
			sysPrompt = fmt.Sprintf(promptTemplate[lang(ctx)], "GPT-3.5")
		}
	} else {
		sysPrompt = fmt.Sprintf(promptTemplate[lang(ctx)], "GPT-4")
	}

	if ctx.Value("user_status") != StatusFree && !strings.Contains(prompt, utils.Delim) {
		var query string
		if len(messages) == 2 && !strings.Contains(messages[0].Content.(string), utils.Delim) {
			query = messages[0].Content.(string) + "\n\n" + prompt
		} else {
			query = prompt
		}
		if ctx.Value("model") == models.GPT3 && lang(ctx) == "uz" {
			ctx = context.WithValue(ctx, "language_code", "en")
			sysPrompt += fmt.Sprintf(searchTemplate[lang(ctx)], r.service.Search(ctx, query))
			ctx = context.WithValue(ctx, "language_code", "uz")
		} else {
			sysPrompt += fmt.Sprintf(searchTemplate[lang(ctx)], r.service.Search(ctx, query))
		}
	}

	messages = append([]types.Message{{Role: "system", Content: sysPrompt}}, messages...)
	messages = append(messages, types.Message{Role: "user", Content: prompt})

	for i := range messages {
		URL, text, found := strings.Cut(messages[i].Content.(string), utils.Delim)
		if found {
			var content []types.Content
			if len(text) > 0 {
				content = append(content, types.Content{Type: "text", Text: text})
			}
			content = append(content, types.Content{Type: "image_url", ImageURL: &types.ImageURL{URL: URL}})
			messages[i].Content = content
			output = context.WithValue(ctx, "model", models.GPT4V)
		}
	}

	if output == nil {
		output = ctx
	}
	return
}

func (r *Redis) StoreHistory(ctx context.Context, prompt, completion string) {
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

func (r *Redis) DeleteHistory(ctx context.Context) {
	r.client.Del(ctx, "context:"+id(ctx))
}
