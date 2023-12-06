package text

var Settings = map[string]string{
	"uz": `Tarif: *bepul*
Versiya: *GPT-3.5*
So'rovlar soni: *%s*
Tugash muddati: *%s*
Rasmni tushunish: *mavjud emas* 🚫
Googledan qidirish: *mavjud emas* 🚫
Fayllar bilan ishlash: *mavjud emas* 🚫
Bot xotirasi: *mavjud emas* 🚫

*Qo'shimcha funksiyalarni* sinab ko'rmoqchimisiz? Unda *pullik tariflarga* o'ting!`,

	"ru": `Подписка: *бесплатная*
Версия: *GPT-3.5*
Количество запросов: *%s*
Срок истечения: *%s*
Понимание изображений: *недоступно* 🚫
Google поиск: *недоступно* 🚫
Работа с файлами: *недоступно* 🚫
Память бота: *недоступно* 🚫

Хотите попробовать *дополнительные функции*? Тогда перейдите на *платные подписки*!`,

	"en": `Subscription: *free*
Version: *GPT-3.5*
Number of requests: *%s*
Expiration date: *%s*
Image understanding: *unavailable* 🚫
Google search: *unavailable* 🚫
Working with files: *unavailable* 🚫
Bot memory: *unavailable* 🚫

Want to try *additional features*? Then switch to the *paid subscriptions*!`,
}

var Settings1 = map[string]string{
	"uz": `Tarif: *cheksiz* ⭐️
Versiya: *GPT-3.5*
So'rovlar soni: *cheklanmagan*
Tugash muddati: *%s*
Rasmni tushunish (faqat rasmdagi matnlarni): *mavjud* ✅
Googleda qidirish: *mavjud* ✅
Fayllar bilan ishlash: *mavjud* ✅
Bot xotirasi: *mavjud* ✅`,

	"ru": `Подписка: *неограниченная* ⭐️
Версия: *GPT-3.5*
Количество запросов: *не ограничено*
Срок истечения: *%s*
Понимание изображений (только текст в изображениях): *доступно* ✅
Google поиск: *доступно* ✅
Работа с файлами: *доступно* ✅
Память бота: *доступно* ✅`,

	"en": `Subscription: *unlimited* ⭐️
Version: *GPT-3.5*
Number of requests: *unlimited*
Expiration date: *%s*
Image understanding (only text in the images): *available* ✅
Google search: *available* ✅
Working with files: *available* ✅
Bot memory: *available* ✅`,
}

var Settings2 = map[string]string{
	"uz": `Tarif: *premium* 🔥
Versiya: *GPT-4 Vision*
So'rovlar soni: *cheklanmagan*
Tugash muddati: *%s*
Rasmni tushunish: *mavjud* ✅
Googleda qidirish: *mavjud* ✅
Fayllar bilan ishlash: *mavjud* ✅
Bot xotirasi: *mavjud* ✅`,

	"ru": `Подписка: *премиум* 🔥
Версия: *GPT-4 Vision*
Количество запросов: *не ограничено*
Срок истечения: *%s*
Понимание изображений: *доступно* ✅
Google поиск: *доступно* ✅
Работа с файлами: *доступно* ✅
Память бота: *доступно* ✅`,

	"en": `Subscription: *premium* 🔥
Version: *GPT-4 Vision*
Number of requests: *unlimited*
Expiration date: *%s*
Image understanding: *available* ✅
Google search: *available* ✅
Working with files: *available* ✅
Bot memory: *available* ✅`,
}
