package text

var Notify1 = map[string]string{
	"uz": `Assalomu alaykum, %s! Bugun sizga qanday yordam bera olishim mumkin?`,
	"ru": `Здравствуйте, %s! Чем я могу вам помочь сегодня?`,
	"en": `Hello, %s! How can I help you today?`,
}

var Notify2 = map[string]string{
	"uz": `*Qadrli %s!* Sizning *%s* tarifingiz *%s*da tugaydi. Yordamchidan to'liq foydalanishni davom ettirish uchun, iltimos, to'lovni amalga oshiring.`,
	"ru": `*Уважаемый(ая) %s!* Ваша *s* подписка закончится в *%s*. Чтобы продолжить пользоваться Yordamchi в полной мере, пожалуйста, совершите оплату.`,
	"en": `*Dear %s!* Your *%s* subscription ends at *%s*. In order to continue using Yordamchi to its fullest, please make a payment.`,
}

var User = map[string]string{
	"uz": `Foydalanuvchi`,
	"ru": `Пользователь`,
	"en": `User`,
}

var UnlimitedSubscription = map[string]string{
	"uz": `cheksiz`,
	"ru": `неограниченная`,
	"en": `unlimited`,
}

var PremiumSubscription = map[string]string{
	"uz": `premium`,
	"ru": `премиум`,
	"en": `premium`,
}
