package text

var Settings = map[string]string{

	"uz": `Tarif: %s
So'rovlar soni: %s
Tugash muddati: %s

Cheksiz so'rovlar kerakmi? Unda premium foydalanuvchi bo'ling!`,

	"ru": `Тариф: %s
Количество запросов: %s
Срок истечения: %s

Нужны безлимитные запросы? Тогда станьте премиум-пользователем!`,

	"en": `Tariff: %s
Number of requests: %s
Expiration date: %s

Need unlimited requests? Then become a premium user!`,
}

var FreeTariff = map[string]string{
	"uz": "oddiy",
	"ru": "обычный",
	"en": "free",
}

var PremiumTariff = map[string]string{
	"uz": "premium 💎",
	"ru": "премиум 💎",
	"en": "premium 💎",
}

var Unlimited = map[string]string{
	"uz": "cheklanmagan",
	"ru": "не ограничено",
	"en": "unlimited",
}
