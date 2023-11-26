package text

var Settings1 = map[string]string{
	"uz": `**Tarif**: bepul
**Versiya**: GPT-3.5
**So'rovlar soni**: %s
**Tugash muddati**: %s
**Rasmni tushunish**: mavjud emas 🚫
**Googleda qidirish**: mavjud emas 🚫

Qo'shimcha funksiyalarni sinab ko'rmoqchimisiz? Unda premium foydalanuvchi bo'ling!`,

	"ru": `**Подписка**: бесплатная
**Версия**: GPT-3.5
**Количество запросов**: %s
**Срок истечения**: %s
**Понимание изображений**: недоступно 🚫
**Google поиск**: недоступно 🚫

Хотите попробовать дополнительные функции? Тогда станьте премиум-пользователем!`,

	"en": `**Subscription**: free
**Version**: GPT-3.5
**Number of requests**: %s
**Expiration date**: %s
**Image understanding**: unavailable 🚫
**Google search**: unavailable 🚫

Want to try additional features? Then become a premium user!`,
}

var Settings2 = map[string]string{
	"uz": `**Tarif**: premium ⭐️
**Versiya**: GPT-4 Vision
**So'rovlar soni**: cheklanmagan
**Tugash muddati**: %s
**Rasmni tushunish**: mavjud ✅
**Googleda qidirish**: mavjud ✅`,

	"ru": `**Подписка**: премиум ⭐️
**Версия**: GPT-4 Vision
**Количество запросов**: не ограничено
**Срок истечения**: %s
**Понимание изображений**: доступно ✅
**Google поиск**: доступно ✅`,

	"en": `**Subscription**: premium ⭐️
**Version**: GPT-4 Vision
**Number of requests**: unlimited
**Expiration date**: %s
**Image understanding**: available ✅
**Google search**: available ✅`,
}
