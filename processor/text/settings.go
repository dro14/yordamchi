package text

var Settings1 = map[string]string{
	"uz": `Tarif: bepul
Versiya: GPT-3.5
Suratni tushunish: mavjud emas 🚫
So'rovlar soni: %s
Tugash muddati: %s

Suratni tushunadigan GPT-4 ni ishlatib ko'rmoqchimisiz? Unda premium foydalanuvchi bo'ling!`,

	"ru": `Подписка: бесплатная
Версия: GPT-3.5
Понимание изображений: недоступно 🚫
Количество запросов: %s
Срок истечения: %s

Хотите попробовать GPT-4, который понимает изображения? Тогда станьте премиум-пользователем!`,

	"en": `Subscription: free
Version: GPT-3.5
Image understanding: unavailable 🚫
Number of requests: %s
Expiration date: %s

Want to try GPT-4 that understands images? Then become a premium user!`,
}

var Settings2 = map[string]string{
	"uz": `Tarif: premium ⭐️
Versiya: GPT-4 & GPT-4 Vision
Suratni tushunish: mavjud ✅
So'rovlar soni: cheklanmagan
Tugash muddati: %s`,

	"ru": `Подписка: премиум ⭐️
Версия: GPT-4 & GPT-4 Vision
Понимание изображений: доступно ✅
Количество запросов: не ограничено
Срок истечения: %s`,

	"en": `Subscription: premium ⭐️
Version: GPT-4 & GPT-4 Vision
Image understanding: available ✅
Number of requests: unlimited
Expiration date: %s`,
}
