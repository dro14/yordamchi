package text

var Settings1 = map[string]string{
	"uz": `Tarif: bepul
Versiya: GPT-3.5
So'rovlar soni: %s
Tugash muddati: %s
Rasmni tushunish: mavjud emas 🚫
Googleda qidirish: mavjud emas 🚫

Qo'shimcha funksiyalarni sinab ko'rmoqchimisiz? Unda premium foydalanuvchi bo'ling!`,

	"ru": `Подписка: бесплатная
Версия: GPT-3.5
Количество запросов: %s
Срок истечения: %s
Понимание изображений: недоступно 🚫
Google поиск: недоступно 🚫

Хотите попробовать дополнительные функции? Тогда станьте премиум-пользователем!`,

	"en": `Subscription: free
Version: GPT-3.5
Number of requests: %s
Expiration date: %s
Image understanding: unavailable 🚫
Google search: unavailable 🚫

Want to try additional features? Then become a premium user!`,
}

var Settings2 = map[string]string{
	"uz": `Tarif: premium ⭐️
Versiya: GPT-4 Vision
Rasmni tushunish: mavjud ✅
So'rovlar soni: cheklanmagan
Tugash muddati: %s`,

	"ru": `Подписка: премиум ⭐️
Версия: GPT-4 Vision
Понимание изображений: доступно ✅
Количество запросов: не ограничено
Срок истечения: %s`,

	"en": `Subscription: premium ⭐️
Version: GPT-4 Vision
Image understanding: available ✅
Number of requests: unlimited
Expiration date: %s`,
}
