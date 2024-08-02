package text

var Memory = map[string]string{
	"uz": `%s

• Botning xotirasiga joylash uchun xabarni quyidagi formatda kiriting:
/system *shu xabar botning xotirasida uzoq muddatga saqlanadi*

• *Botning xotirasini* tozalash uchun «💬 *Yangi* 💬» tugmasini bosing`,

	"ru": `%s

• Для сохранения сообщения в памяти бота введите его в следующем формате:
/system *это сообщение будет сохранено в памяти бота на длительный срок*

• Чтобы очистить *память бота* нажмите кнопку «💬 *Новый* 💬»`,

	"en": `%s

• To save a message in the bot memory, enter it in the following format:
/system *this message will be saved in the bot memory for a long time*

• To clear the *bot memory*, click the «💬 *New* 💬» button`,
}

var MemorySystem = map[string]string{
	"uz": "Hozir botning xotirasida quyidagi xabar saqlanyapti: *%s*",
	"ru": "Сейчас в памяти бота сохранено следующее сообщение: *%s*",
	"en": "Currently bot memory holds the following message: *%s*",
}

var MemoryEmpty = map[string]string{
	"uz": "*Hozir botning xotirasi bo'sh*",
	"ru": "*Сейчас память бота пуста*",
	"en": "*Currently bot memory is empty*",
}
