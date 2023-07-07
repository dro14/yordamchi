package button

import (
	"context"

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
	if redis.Model(ctx) == "gpt-3.5-turbo" {
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
			Data: []byte("gpt-3.5-turbo"),
		},
	)
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonCallback{
			Text: texts[1],
			Data: []byte("gpt-4"),
		},
	)

	keyboard := &tg.ReplyInlineMarkup{}
	keyboard.Rows = append(keyboard.Rows, row)
	return keyboard
}

func Translate(lang string) *tg.ReplyInlineMarkup {

	keyboard := &tg.ReplyInlineMarkup{}

	text := map[string]string{
		"uz": "✅ Yoqish ✅",
		"ru": "✅ Включить ✅",
		"en": "✅ Enable ✅",
	}
	row := tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonCallback{
			Text: text[lang],
			Data: []byte("enable"),
		},
	)
	keyboard.Rows = append(keyboard.Rows, row)

	text = map[string]string{
		"uz": "🚫 O'chirish 🚫",
		"ru": "🚫 Выключить 🚫",
		"en": "🚫 Disable 🚫",
	}
	row = tg.KeyboardButtonRow{}
	row.Buttons = append(row.Buttons,
		&tg.KeyboardButtonCallback{
			Text: text[lang],
			Data: []byte("disable"),
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
