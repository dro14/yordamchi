package button

import (
	"github.com/dro14/yordamchi/lib/constants"
	"github.com/gotd/td/tg"
)

func NewChat(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "💬 Yangi suhbat 💬",
		"ru": "💬 Новый разговор 💬",
		"en": "💬 New chat 💬",
	}
	return DataButton(text[lang], "new_chat")
}

func Start(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "❔ Qanday ishlatish ❔",
		"ru": "❔ Как пользоваться ❔",
		"en": "❔ How to use ❔",
	}
	return DataButton(text[lang], "examples")
}

func Examples(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "📝 Bot haqida ma'lumot 📝",
		"ru": "📝 Информация о боте 📝",
		"en": "📝 Information about the bot 📝",
	}
	return DataButton(text[lang], "help")
}

func Settings(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "💎 So'rovlar sotib olish 💎",
		"ru": "💎 Купить запросы 💎",
		"en": "💎 Buy requests 💎",
	}
	return DataButton(text[lang], "premium")
}

func Premium(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "🚀 Tasdiqlash 🚀",
		"ru": "🚀 Подтверждать 🚀",
		"en": "🚀 Confirm 🚀",
	}
	return DataButton(text[lang], "1000000")
}

func Exhausted(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "⭐ Premium bo'lish ⭐",
		"ru": "⭐ Стать премиумом ⭐",
		"en": "⭐ Become premium ⭐",
	}
	return DataButton(text[lang], "premium")
}

func Blocked(lang string) *tg.ReplyInlineMarkup {
	text := map[string]string{
		"uz": "👤 Admin 👤",
		"ru": "👤 Админ 👤",
		"en": "👤 Admin 👤",
	}
	return URLButton(text[lang], constants.AdminURL)
}

func URLButton(text, url string) *tg.ReplyInlineMarkup {

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

func DataButton(text, data string) *tg.ReplyInlineMarkup {

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
