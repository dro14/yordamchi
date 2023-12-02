package text

var Memory = map[string]string{
	"uz": "Botning xotirasida quyidagi %s mavjud:\n\n*%s*",
	"ru": "В памяти бота находится %s:\n\n*%s*",
	"en": "Bot memory contains the following %s:\n\n*%s*",
}

var MemoryEmpty = map[string]string{
	"uz": "Hozir botning xotirasi bo'sh va *Google* qo'shimcha ",
	"ru": "Память бота пуста. Поэтому дополнительная информация берется из *Google*.",
	"en": "Bot memory is empty. Therefore, additional information is taken from *Google*.",
}
