package processor

import (
	"context"

	"github.com/dro14/yordamchi/storage/redis"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) newChatButton(ctx context.Context) *tgbotapi.InlineKeyboardMarkup {
	text := map[string]string{
		"uz": "💬 Yangi suhbat 💬",
		"ru": "💬 Новый разговор 💬",
		"en": "💬 New chat 💬",
	}
	return data("new_chat", text[lang(ctx)])
}

func (p *Processor) startButton(ctx context.Context) *tgbotapi.InlineKeyboardMarkup {
	text := map[string]string{
		"uz": "❔ Qanday ishlatish ❔",
		"ru": "❔ Как пользоваться ❔",
		"en": "❔ How to use ❔",
	}
	return data("examples", text[lang(ctx)])
}

func (p *Processor) settingsButton(ctx context.Context) *tgbotapi.InlineKeyboardMarkup {
	if ctx.Value("user_status") == redis.StatusPremium {
		return nil
	}
	text := map[string]string{
		"uz": "⭐️ Premium tarif ⭐️",
		"ru": "⭐️ Премиум подписка ⭐️",
		"en": "⭐️ Premium subscription ⭐️",
	}
	return data("settings", text[lang(ctx)])
}

func (p *Processor) languageButtons() *tgbotapi.InlineKeyboardMarkup {
	return data("uz", "ru", "en", "🇺🇿 O'zbekcha 🇺🇿", "🇷🇺 Русский 🇷🇺", "🇺🇸 English 🇺🇸")
}

func (p *Processor) examplesButton(ctx context.Context) *tgbotapi.InlineKeyboardMarkup {
	text := map[string]string{
		"uz": "📝 Bot haqida ma'lumot 📝",
		"ru": "📝 Информация о боте 📝",
		"en": "📝 Information about the bot 📝",
	}
	return data("help", text[lang(ctx)])
}

func (p *Processor) premiumButtons(ctx context.Context) *tgbotapi.InlineKeyboardMarkup {
	text := map[string][]string{
		"uz": {"⭐️ Kunlik ⭐️", "🔥 Haftalik 🔥", "🚀 Oylik 🚀"},
		"ru": {"⭐️ Дневная ⭐️", "🔥 Недельная 🔥", "🚀 Месячная 🚀"},
		"en": {"⭐️ Daily ⭐️", "🔥 Weekly 🔥", "🚀 Monthly 🚀"},
	}
	args := make([]string, 6)
	args[0] = p.payme.CheckoutURL(ctx, 1000000, "daily:gpt-4")
	args[1] = p.payme.CheckoutURL(ctx, 5000000, "weekly:gpt-4")
	args[2] = p.payme.CheckoutURL(ctx, 15000000, "monthly:gpt-4")
	args[3] = text[lang(ctx)][0]
	args[4] = text[lang(ctx)][1]
	args[5] = text[lang(ctx)][2]
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
