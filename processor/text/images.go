package text

var Images = map[string]string{
	"uz": `So'zlar rasmga aylanganini ko'rganmisiz? Yo'q bo'lsa, *DALL-E 3* ni kutib oling!

*DALL-E 3* - bu so'zlar yordamida rasm generatsiya qiladigan eng kuchli sun'iy intellektlardan biri. Endi uning yordamida siz yuqoridagidek aql bovar qilmas rasmlar yaratishingiz mumkin. Hattoki storislar uchun yozuvli rasmlar ham, faqat yozuvlari ingliz tilida bo'ladi :) 

Rasm yaratish uchun so'rovingizni quyidagi shaklda yuboring:

/generate *Qumli sokin sohil bo'yida quyosh botyapti. Qayiqda odamlar kechgi ovqatga yig'ilgan*

/generate *Ko'prikda turibsiz. Kechki shaharda mayin yomg'ir yog'yapti. Odamlar ishdan qaytyapti*

/generate *"Tezlik" nomli YouTube kanalga banner*

Rasmni tasvirlashda qancha ko'p detallari bo'lsa, natija ham shuncha yaxshi bo'ladi. **Tabiiy** uslubda rasmlar haqiqiy hayotga yaqinroq bo'ladi, **yorqin** uslubda esa rasmlar ko'proq rang-barang bo'ladi.

Mavjud generatsiyalar soni: *%d*

  20 000 so'm -  *10ta rasm*
  80 000 so'm -  *50ta rasm*
130 000 so'm - *100ta rasm*`,

	"ru": `Вы видели, как слова превращаются в изображение? Если нет, то встречайте *DALL-E 3*!

*DALL-E 3* - один из самых мощных искусственных интеллектов, генерирующий изображения из слов. Теперь с его помощью вы можете создавать ошеломляющие изображения, как те выше. Даже изображения с текстами для сторис, но только тексты будут на английском :)

Для создания изображения отправьте запрос в следующем формате:

/generate *Закат солнца за песчаным берегом. На лодке люди собрались на мероприятие*

/generate *Вы стоите на мосту. Вечерний город затянуло туманом. Люди возвращаются с работы*

/generate *Баннер для YouTube канала "Скорость"*

Чем больше деталей в описании изображения, тем лучше результат. В **натуральном** стиле изображения будут ближе к реальной жизни, а в **ярком** стиле - более красочными.

Количество доступных генераций: *%d*

  20 000 сум -  *10 изображений*
  80 000 сум -  *50 изображений*
130 000 сум - *100 изображений*`,

	"en": `Have you seen how words turn into images? If not, then meet *DALL-E 3*!

*DALL-E 3* is one of the most powerful artificial intelligences that generates images. Now with its help you can create stunning images, even those with texts for stories, but the texts will be only in English :)

To create an image, send a request in the following format:

/generate *Sunset over the sandy beach. People are gathering for an evening meal on a boat. The holiday is in full swing*

/generate *You are standing on a bridge. Evening's city was covered with fog. People are returning from work. In apartments, the lights are on*

/generate *Banner for YouTube channel "Speed"*

The more details in the description of the image, the better the result. In the **natural** style, the images will be closer to real life, while in the **vivid** style - brighter and more colorful.

Number of available generations: *%d*

  20 000 UZS -  *10 images*
  80 000 UZS -  *50 images*
130 000 UZS - *100 images*`,
}

var ImagesPayments = map[string]string{
	"uz": `*DALL-E* bilan tasavvuringizni rasmlarga aylantiring!

  20 000 so'm -  *10ta rasm*
  80 000 so'm -  *50ta rasm*
130 000 so'm - *100ta rasm*`,
	"ru": `Превращайте свое воображение в изображения с *DALL-E*!

  20 000 сум -  *10 изображений*
  80 000 сум -  *50 изображений*
130 000 сум - *100 изображений*`,
	"en": `Turn your imagination into images with *DALL-E*!

  20 000 UZS -  *10 images*
  80 000 UZS -  *50 images*
130 000 UZS - *100 images*`,
}
