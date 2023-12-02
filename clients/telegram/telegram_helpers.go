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
		{Command: "memory", Description: "Botning xotirasi"},
		{Command: "examples", Description: "Foydalanish misollari"},
		{Command: "unlimited", Description: "Cheksiz tarif"},
		{Command: "premium", Description: "Premium tarif"},
		{Command: "images", Description: "Rasm generatsiyasi"},
	},
	"ru": {
		{Command: "start", Description: "Запустить бота"},
		{Command: "help", Description: "Информация о боте"},
		{Command: "settings", Description: "Настройки бота"},
		{Command: "language", Description: "Изменить язык"},
		{Command: "memory", Description: "Память бота"},
		{Command: "examples", Description: "Примеры использования"},
		{Command: "unlimited", Description: "Неограниченная подписка"},
		{Command: "premium", Description: "Премиум подписка"},
		{Command: "images", Description: "Генерация изображений"},
	},
	"en": {
		{Command: "start", Description: "Start the bot"},
		{Command: "help", Description: "Information about the bot"},
		{Command: "settings", Description: "Bot settings"},
		{Command: "language", Description: "Change language"},
		{Command: "memory", Description: "Bot memory"},
		{Command: "examples", Description: "Usage examples"},
		{Command: "unlimited", Description: "Unlimited subscription"},
		{Command: "premium", Description: "Premium subscription"},
		{Command: "images", Description: "Image generation"},
	},
}

func id(ctx context.Context) int64 {
	return ctx.Value("user_id").(int64)
}
