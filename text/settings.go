package text

var Settings1 = map[string]string{

	"uz": `Tarif: %s
So'rovlar soni: %s
Tugash muddati: %s`,

	"ru": `Тариф: %s
Количество запросов: %s
Срок истечения: %s`,

	"en": `Tariff: %s
Number of requests: %s
Expiration date: %s`,
}

var FreeTariff = map[string]string{
	"uz": "oddiy",
	"ru": "обычный",
	"en": "free",
}

var PremiumTariff = map[string]string{
	"uz": "premium",
	"ru": "премиум",
	"en": "premium",
}

var Unlimited = map[string]string{
	"uz": "cheklanmagan",
	"ru": "не ограничено",
	"en": "unlimited",
}

var Settings2 = map[string]string{
	"uz": `Mavjud tokenlar: %d`,
	"ru": `Доступно токенов: %d`,
	"en": `Available tokens: %d`,
}
