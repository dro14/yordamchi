package text

var Translate = map[string]string{
	"uz": `Agar Siz botni o'zbekchada ishlatayotgan bo'lsangiz va u ma'nosiz javoblar qaytarayotgan bo'lsa, tarjimon funksiyasini yoqishingiz mumkin. Shunda bot sifatliroq o'zbekcha javoblar beradi.

Tarjimon yoqilganda, strim ishlamaydi.`,

	"ru": `Если Вы используете бот на узбекском языке и он возвращает бессмысленные ответы, Вы можете включить функцию переводчика. Тогда бот будет отвечать на узбекском более качественно.

При включенном переводчике, стрим не работает.`,

	"en": `If you are using the bot in Uzbek and it is returning meaningless answers, you can turn on the translator function. Then the bot will respond in Uzbek more qualitatively.

When the translator is turned on, the stream does not work.`,
}

var TranslatorEnabled = map[string]string{
	"uz": "Tarjimon yoqildi",
	"ru": "Переводчик включен",
	"en": "Translator enabled",
}

var TranslatorDisabled = map[string]string{
	"uz": "Tarjimon o'chirildi",
	"ru": "Переводчик выключен",
	"en": "Translator disabled",
}
