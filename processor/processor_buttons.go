package processor

import (
	"context"

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
	unlimitedText := map[string]string{
		"uz": "⭐️ Cheksiz obuna ⭐️",
		"ru": "⭐️ Безлимитная подписка ⭐️",
		"en": "⭐️ Unlimited subscription ⭐️",
	}
	premiumText := map[string]string{
		"uz": "🔥 Premium obuna 🔥",
		"ru": "🔥 Премиум подписка 🔥",
		"en": "🔥 Premium subscription 🔥",
	}
	return data("settings1", "settings2", unlimitedText[lang(ctx)], premiumText[lang(ctx)])
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

func (p *Processor) unlimitedButtons(ctx context.Context) *tgbotapi.InlineKeyboardMarkup {
	text := map[string][]string{
		"uz": {"⭐️ Haftalik ⭐️", "🔥 Oylik 🔥"},
		"ru": {"⭐️ Недельная ⭐️", "🔥 Месячная 🔥"},
		"en": {"⭐️ Weekly ⭐️", "🔥 Monthly 🔥"},
	}
	args := make([]string, 4)
	args[0] = p.payme.CheckoutURL(ctx, 1000000, "weekly:unlimited")
	args[1] = p.payme.CheckoutURL(ctx, 3000000, "monthly:unlimited")
	args[2] = text[lang(ctx)][0]
	args[3] = text[lang(ctx)][1]
	return url(args...)
}

func (p *Processor) premiumButtons(ctx context.Context) *tgbotapi.InlineKeyboardMarkup {
	text := map[string][]string{
		"uz": {"⭐️ Kunlik ⭐️", "🔥 Haftalik 🔥", "🚀 Oylik 🚀"},
		"ru": {"⭐️ Дневная ⭐️", "🔥 Недельная 🔥", "🚀 Месячная 🚀"},
		"en": {"⭐️ Daily ⭐️", "🔥 Weekly 🔥", "🚀 Monthly 🚀"},
	}
	args := make([]string, 6)
	args[0] = p.payme.CheckoutURL(ctx, 1000000, "daily:premium")
	args[1] = p.payme.CheckoutURL(ctx, 5000000, "weekly:premium")
	args[2] = p.payme.CheckoutURL(ctx, 15000000, "monthly:premium")
	args[3] = text[lang(ctx)][0]
	args[4] = text[lang(ctx)][1]
	args[5] = text[lang(ctx)][2]
	return url(args...)
}

func (p *Processor) imageButtons(ctx context.Context) *tgbotapi.InlineKeyboardMarkup {
	text := map[string][]string{
		"uz": {"⭐️ 10ta rasm ⭐️", "🔥 50ta rasm 🔥", "🚀 100ta rasm 🚀"},
		"ru": {"⭐️ 10 изображений ⭐️", "🔥 50 изображений 🔥", "🚀 100 изображений 🚀"},
		"en": {"⭐️ 10 images ⭐️", "🔥 50 images 🔥", "🚀 100 images 🚀"},
	}
	args := make([]string, 6)
	args[0] = p.payme.CheckoutURL(ctx, 2000000, "10:images")
	args[1] = p.payme.CheckoutURL(ctx, 8000000, "50:images")
	args[2] = p.payme.CheckoutURL(ctx, 13000000, "100:images")
	args[3] = text[lang(ctx)][0]
	args[4] = text[lang(ctx)][1]
	args[5] = text[lang(ctx)][2]
	return url(args...)
}

func (p *Processor) generateButtons(ctx context.Context) *tgbotapi.InlineKeyboardMarkup {
	vividText := map[string]string{
		"uz": "yorqin",
		"ru": "яркого",
		"en": "vivid",
	}
	naturalText := map[string]string{
		"uz": "tabiiy",
		"ru": "натурального",
		"en": "natural",
	}
	vivid := "vivid"
	natural := "natural"
	row := [][]tgbotapi.InlineKeyboardButton{{
		{Text: vividText[lang(ctx)], CallbackData: &vivid},
		{Text: naturalText[lang(ctx)], CallbackData: &natural},
	}}
	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: row}
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
