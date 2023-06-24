package button

import (
	"context"

	"github.com/dro14/yordamchi/payme"
	"github.com/gotd/td/tg"
)

func NewChat(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "💬 Yangi suhbat 💬",
		"ru": "💬 Новый разговор 💬",
		"en": "💬 New chat 💬",
	}
	return data(text[lang], "new_chat")
}

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

func Settings(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "💎 Premium bo'lish 💎",
		"ru": "💎 Стать премиумом 💎",
		"en": "💎 Become premium 💎",
	}
	return data(text[lang], "premium")
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

func Donate(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "😇 Xayriya 😇",
		"ru": "😇 Донат 😇",
		"en": "😇 Donate 😇",
	}
	return url(text[lang], "https://payme.uz/60d6dbeb3632e1ceb8664de3")
}

func Blocked(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "👤 Admin 👤",
		"ru": "👤 Админ 👤",
		"en": "👤 Admin 👤",
	}
	return url(text[lang], "https://t.me/yordamchiga_yordam")
}

func url(text, url string) *tg.ReplyInlineMarkup {

	row := tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonURL{
			Text: text,
			URL:  url,
		},
	)

	keyboard := &tg.ReplyInlineMarkup{}
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
