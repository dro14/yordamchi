package button

import (
	"context"

	"github.com/dro14/yordamchi/lib/models"
	"github.com/dro14/yordamchi/payme"
	"github.com/dro14/yordamchi/redis"
	"github.com/gotd/td/tg"
)

func Start(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "❔ Qanday ishlatish ❔",
		"ru": "❔ Как пользоваться ❔",
		"en": "❔ How to use ❔",
	}
	return data(text[lang], "examples")
}

func Examples(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "📝 Bot haqida ma'lumot 📝",
		"ru": "📝 Информация о боте 📝",
		"en": "📝 Information about the bot 📝",
	}
	return data(text[lang], "help")
}

func Settings(ctx context.Context) *tg.ReplyInlineMarkup {

	texts := make([]string, 2)
	if redis.Model(ctx) == models.GPT3 {
		texts[0] = "GPT-3.5 ✅"
		texts[1] = "GPT-4"
	} else {
		texts[0] = "GPT-3.5"
		texts[1] = "GPT-4 ✅"
	}

	row := tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonCallback{
			Text: texts[0],
			Data: []byte(models.GPT4),
		},
	)
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonCallback{
			Text: texts[1],
			Data: []byte(models.GPT4),
		},
	)

	keyboard := &tg.ReplyInlineMarkup{}
	keyboard.Rows = append(keyboard.Rows, row)
	return keyboard
}

func Language() *tg.ReplyInlineMarkup {

	keyboard := &tg.ReplyInlineMarkup{}
	row := tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonCallback{
			Text: "🇺🇿 O'zbekcha 🇺🇿",
			Data: []byte("uz"),
		},
	)
	keyboard.Rows = append(keyboard.Rows, row)

	row = tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonCallback{
			Text: "🇷🇺 Русский 🇷🇺",
			Data: []byte("ru"),
		},
	)
	keyboard.Rows = append(keyboard.Rows, row)

	row = tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonCallback{
			Text: "🇺🇸 English 🇺🇸",
			Data: []byte("en"),
		},
	)
	keyboard.Rows = append(keyboard.Rows, row)
	return keyboard
}

func Premium(ctx context.Context, lang string) *tg.ReplyInlineMarkup {

	keyboard := &tg.ReplyInlineMarkup{}

	var weekly = map[string]string{
		"uz": "⭐ Haftalik ⭐",
		"ru": "⭐ Недельный ⭐",
		"en": "⭐ Weekly ⭐",
	}
	row := tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonURL{
			Text: weekly[lang],
			URL:  payme.CheckoutURL(ctx, 1000000, "weekly"),
		},
	)
	keyboard.Rows = append(keyboard.Rows, row)

	var monthly = map[string]string{
		"uz": "🔥 Oylik 🔥",
		"ru": "🔥 Месячный 🔥",
		"en": "🔥 Monthly 🔥",
	}
	row = tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonURL{
			Text: monthly[lang],
			URL:  payme.CheckoutURL(ctx, 3000000, "monthly"),
		},
	)
	keyboard.Rows = append(keyboard.Rows, row)

	return keyboard
}

func GPT4(ctx context.Context, lang string) *tg.ReplyInlineMarkup {

	keyboard := &tg.ReplyInlineMarkup{}

	var ten = map[string]string{
		"uz": "⭐ 10,000 ta token ⭐",
		"ru": "⭐ 10,000 токенов ⭐",
		"en": "⭐ 10,000 tokens ⭐",
	}
	row := tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonURL{
			Text: ten[lang],
			URL:  payme.CheckoutURL(ctx, 1000000, "gpt-4"),
		},
	)
	keyboard.Rows = append(keyboard.Rows, row)

	var thirty = map[string]string{
		"uz": "🔥 30,000 ta token 🔥",
		"ru": "🔥 30,000 токенов 🔥",
		"en": "🔥 30,000 tokens 🔥",
	}
	row = tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonURL{
			Text: thirty[lang],
			URL:  payme.CheckoutURL(ctx, 3000000, "gpt-4"),
		},
	)
	keyboard.Rows = append(keyboard.Rows, row)

	var hundred = map[string]string{
		"uz": "🚀 100,000 ta token 🚀",
		"ru": "🚀 100,000 токенов 🚀",
		"en": "🚀 100,000 tokens 🚀",
	}
	row = tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonURL{
			Text: hundred[lang],
			URL:  payme.CheckoutURL(ctx, 10000000, "gpt-4"),
		},
	)
	keyboard.Rows = append(keyboard.Rows, row)

	return keyboard
}

func data(text, data string) *tg.ReplyInlineMarkup {

	row := tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonCallback{
			Text: text,
			Data: []byte(data),
		},
	)

	keyboard := &tg.ReplyInlineMarkup{}
	keyboard.Rows = append(keyboard.Rows, row)
	return keyboard
}
