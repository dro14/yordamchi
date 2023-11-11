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
	return data("examples", text[lang])
}

func Examples(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "📝 Bot haqida ma'lumot 📝",
		"ru": "📝 Информация о боте 📝",
		"en": "📝 Information about the bot 📝",
	}
	return data("help", text[lang])
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
			Data: []byte(models.GPT3),
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
	return data("uz", "ru", "en", "🇺🇿 O'zbekcha 🇺🇿", "🇷🇺 Русский 🇷🇺", "🇺🇸 English 🇺🇸")
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

func data(args ...string) *tg.ReplyInlineMarkup {
	keyboard := &tg.ReplyInlineMarkup{}
	n := len(args) / 2
	for i := 0; i < n; i++ {
		row := tg.KeyboardButtonRow{}
		row.Buttons = append(row.Buttons,
			&tg.KeyboardButtonCallback{
				Text: args[i+n],
				Data: []byte(args[i]),
			},
		)
		keyboard.Rows = append(keyboard.Rows, row)
	}
	return keyboard
}

//func Premium(ctx context.Context, lang string) *tg.ReplyInlineMarkup {
//	text := map[string][]string{
//		"uz": {"⭐ Kunlik ⭐", "🔥 Haftalik 🔥", "🚀 Oylik 🚀"},
//		"ru": {"⭐ Суточный ⭐", "🔥 Недельный 🔥", "🚀 Месячный 🚀"},
//		"en": {"⭐ Daily ⭐", "🔥 Weekly 🔥", "🚀 Monthly 🚀"},
//	}
//	args := make([]string, 6)
//	args[0] = payme.CheckoutURL(ctx, 599000, "daily")
//	args[1] = payme.CheckoutURL(ctx, 1999000, "weekly")
//	args[2] = payme.CheckoutURL(ctx, 5999000, "monthly")
//	args[3] = text[lang][0]
//	args[4] = text[lang][1]
//	args[5] = text[lang][2]
//	return url(args...)
//}
//
//func url(args ...string) *tg.ReplyInlineMarkup {
//	keyboard := &tg.ReplyInlineMarkup{}
//	n := len(args) / 2
//	for i := 0; i < n; i++ {
//		row := tg.KeyboardButtonRow{}
//		row.Buttons = append(row.Buttons,
//			&tg.KeyboardButtonURL{
//				Text: args[i+n],
//				URL:  args[i],
//			},
//		)
//		keyboard.Rows = append(keyboard.Rows, row)
//	}
//	return keyboard
//}
//
