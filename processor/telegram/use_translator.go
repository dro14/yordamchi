package telegram

import (
	"context"
	"fmt"
	"github.com/dro14/yordamchi/lib/models"

	"github.com/dro14/yordamchi/client/translator"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/processor/openai"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func UseTranslator(ctx context.Context, message *tgbotapi.Message, stats *types.Stats) []string {

	prompt, err := translator.Translate("auto", "en", message.Text)
	if err != nil {
		return []string{text.TranslationFailed[lang(ctx)]}
	}

	messages := redis.LoadContext(ctx, prompt)
	completion, err := openai.Process(ctx, messages, stats)
	if err != nil {
		return []string{completion}
	}

	translation, err := translator.Translate("auto", lang(ctx), completion)
	if err != nil {
		return []string{text.TranslationFailed[lang(ctx)]}
	}

	tokensUsed := stats.PromptTokens + stats.CompletionTokens
	if ctx.Value("model") == models.GPT4 {
		translation = fmt.Sprintf(text.TokensUsed[lang(ctx)], translation, tokensUsed)
	}

	redis.Decrement(ctx, tokensUsed)
	redis.StoreContext(ctx, prompt, completion)
	return functions.Slice(translation)
}
