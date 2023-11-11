package button

import (
	"context"

	"github.com/dro14/yordamchi/lib/models"
	"github.com/dro14/yordamchi/payme"
	"github.com/dro14/yordamchi/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewChat(lang string) *tgbotapi.InlineKeyboardMarkup {
	text := map[string]string{
		"uz": "💬 Yangi suhbat 💬",
		"ru": "💬 Новый разговор 💬",
		"en": "💬 New chat 💬",
	}
	return data("new_chat", text[lang])
}

func Start(lang string) *tgbotapi.InlineKeyboardMarkup {
	text := map[string]string{
		"uz": "❔ Qanday ishlatish ❔",
		"ru": "❔ Как пользоваться ❔",
		"en": "❔ How to use ❔",
	}
	return data("examples", text[lang])
}

func Settings(ctx context.Context) *tgbotapi.InlineKeyboardMarkup {
	gpt3 := models.GPT3
	gpt4 := models.GPT4
	text := make(map[string]string)
	if redis.Model(ctx) == models.GPT3 {
		text[gpt3] = "GPT-3.5 ✅"
		text[gpt4] = "GPT-4"
	} else {
		text[gpt3] = "GPT-3.5"
		text[gpt4] = "GPT-4 ✅"
	}
	keyboard := [][]tgbotapi.InlineKeyboardButton{{
		{Text: text[gpt3], CallbackData: &gpt3},
		{Text: text[gpt4], CallbackData: &gpt4},
	}}
	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: keyboard}
}

func Language() *tgbotapi.InlineKeyboardMarkup {
	return data("uz", "ru", "en", "🇺🇿 O'zbekcha 🇺🇿", "🇷🇺 Русский 🇷🇺", "🇺🇸 English 🇺🇸")
}

func Examples(lang string) *tgbotapi.InlineKeyboardMarkup {
	text := map[string]string{
		"uz": "📝 Bot haqida ma'lumot 📝",
		"ru": "📝 Информация о боте 📝",
		"en": "📝 Information about the bot 📝",
	}
	return data("help", text[lang])
}

func Premium(ctx context.Context, lang string) *tgbotapi.InlineKeyboardMarkup {
	text := map[string][]string{
		"uz": {"⭐ Haftalik ⭐", "🔥 Oylik 🔥"},
		"ru": {"⭐ Недельный ⭐", "🔥 Месячный 🔥"},
		"en": {"⭐ Weekly ⭐", "🔥 Monthly 🔥"},
	}
	args := make([]string, 4)
	args[0] = payme.CheckoutURL(ctx, 1000000, "weekly")
	args[1] = payme.CheckoutURL(ctx, 3000000, "monthly")
	args[2] = text[lang][0]
	args[3] = text[lang][1]
	return url(args...)
}

func GPT4(ctx context.Context, lang string) *tgbotapi.InlineKeyboardMarkup {
	text := map[string][]string{
		"uz": {"⭐ 10,000 ta token ⭐", "🔥 30,000 ta token 🔥", "🚀 100,000 ta token 🚀"},
		"ru": {"⭐ 10,000 токенов ⭐", "🔥 30,000 токенов 🔥", "🚀 100,000 токенов 🚀"},
		"en": {"⭐ 10,000 tokens ⭐", "🔥 30,000 tokens 🔥", "🚀 100,000 tokens 🚀"},
	}
	args := make([]string, 6)
	args[0] = payme.CheckoutURL(ctx, 1000000, "gpt-4")
	args[1] = payme.CheckoutURL(ctx, 3000000, "gpt-4")
	args[2] = payme.CheckoutURL(ctx, 10000000, "gpt-4")
	args[3] = text[lang][0]
	args[4] = text[lang][1]
	args[5] = text[lang][2]
	return url(args...)
}

func data(args ...string) *tgbotapi.InlineKeyboardMarkup {
	keyboard := &tgbotapi.InlineKeyboardMarkup{}
	n := len(args) / 2
	for i := 0; i < n; i++ {
		row := []tgbotapi.InlineKeyboardButton{
			{Text: args[i+n], CallbackData: &args[i]},
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}
	return keyboard
}

func url(args ...string) *tgbotapi.InlineKeyboardMarkup {
	keyboard := &tgbotapi.InlineKeyboardMarkup{}
	n := len(args) / 2
	for i := 0; i < n; i++ {
		row := []tgbotapi.InlineKeyboardButton{
			{Text: args[i+n], URL: &args[i]},
		}
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}
	return keyboard
}
