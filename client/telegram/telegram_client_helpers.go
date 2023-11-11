package telegram

import "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var commands = map[string][]tgbotapi.BotCommand{
	"uz": {
		{Command: "start", Description: "Botni ishga tushirish"},
		{Command: "help", Description: "Bot haqida ma'lumot"},
		{Command: "settings", Description: "Botning sozlamalari"},
		{Command: "language", Description: "Tilni o'zgartirish"},
		{Command: "examples", Description: "Foydalanish misollari"},
		{Command: "premium", Description: "Cheksiz so'rovlar"},
		{Command: "gpt4", Description: "Eng kuchlisi"},
		{Command: "image", Description: "Rasm generatsiyasi"},
	},
	"ru": {
		{Command: "start", Description: "Запустить бота"},
		{Command: "help", Description: "Информация о боте"},
		{Command: "settings", Description: "Настройки бота"},
		{Command: "language", Description: "Изменить язык"},
		{Command: "examples", Description: "Примеры использования"},
		{Command: "premium", Description: "Безлимитные запросы"},
		{Command: "gpt4", Description: "Самый мощный"},
		{Command: "image", Description: "Генерация изображений"},
	},
	"en": {
		{Command: "start", Description: "Start the bot"},
		{Command: "help", Description: "Information about the bot"},
		{Command: "settings", Description: "Bot settings"},
		{Command: "language", Description: "Change language"},
		{Command: "examples", Description: "Usage examples"},
		{Command: "premium", Description: "Unlimited requests"},
		{Command: "gpt4", Description: "The most powerful"},
		{Command: "image", Description: "Image generation"},
	},
}

func newChatButton(lang string) *tgbotapi.InlineKeyboardMarkup {
	text := map[string]string{
		"uz": "💬 Yangi suhbat 💬",
		"ru": "💬 Новый разговор 💬",
		"en": "💬 New chat 💬",
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				text[lang], "new_chat",
			),
		),
	)
	return &keyboard
}
