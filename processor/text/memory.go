package text

var Memory = map[string]string{
	"uz": `%s

Qo'shimcha ma'lumotlar manbasi: *%s*

• Botning xotirasiga joylash uchun xabarni quyidagi formatda kiriting:
/system *shu xabar botning xotirasida uzoq muddatga saqlanadi*

• *Qo'shimcha ma'lumotlar manbasini* kiritish uchun shunchaki botga fayl yuboring, *20 MB*dan oshmagan

• *Botning xotirasi* va *qo'shimcha ma'lumotlar manbasini* tozalash uchun «💬 *Yangi* 💬» tugmasini bosing`,

	"ru": `%s

Источник дополнительной информации: *%s*

• Для сохранения сообщения в памяти бота введите его в следующем формате:
/system *это сообщение будет сохранено в памяти бота на длительный срок*

• Для добавления *источника дополнительной информации* просто отправьте боту файл, не более чем *20 МБ*

• Чтобы очистить *память бота* и *источник дополнительной информации* нажмите кнопку «💬 *Новый* 💬»`,

	"en": `%s

Source of additional information: *%s*

• To save a message in the bot memory, enter it in the following format:
/system *this message will be saved in the bot memory for a long time*

• To add a *source of additional information*, just send a file to the bot, no more than *20 MB*

• To clear the *bot memory* and the *source of additional information*, click the «💬 *New* 💬» button`,
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
