package telegram

import (
	"context"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var commands = map[string][]tgbotapi.BotCommand{
	"uz": {
		{Command: "start", Description: "Botni ishga tushirish"},
		{Command: "help", Description: "Bot haqida ma'lumot"},
		{Command: "settings", Description: "Botning sozlamalari"},
		{Command: "language", Description: "Tilni o'zgartirish"},
		{Command: "examples", Description: "Foydalanish misollari"},
		{Command: "premium", Description: "Qo'shimcha funksiyalar"},
		{Command: "image", Description: "Tasvir generatsiyasi"},
	},
	"ru": {
		{Command: "start", Description: "Запустить бота"},
		{Command: "help", Description: "Информация о боте"},
		{Command: "settings", Description: "Настройки бота"},
		{Command: "language", Description: "Изменить язык"},
		{Command: "examples", Description: "Примеры использования"},
		{Command: "premium", Description: "Дополнительные функции"},
		{Command: "image", Description: "Генерация изображений"},
	},
	"en": {
		{Command: "start", Description: "Start the bot"},
		{Command: "help", Description: "Information about the bot"},
		{Command: "settings", Description: "Bot settings"},
		{Command: "language", Description: "Change language"},
		{Command: "examples", Description: "Usage examples"},
		{Command: "premium", Description: "Additional features"},
		{Command: "image", Description: "Image generation"},
	},
}

func id(ctx context.Context) int64 {
	return ctx.Value("user_id").(int64)
}
