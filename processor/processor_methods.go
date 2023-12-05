package processor

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/dro14/yordamchi/processor/text"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Processor) exhausted(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.Exhausted[lang(ctx)], 0, p.settingsButton(ctx))
	if err != nil {
		log.Println("can't send exhausted message")
	}
}

func (p *Processor) paidFeature(ctx context.Context) {
	_, err := p.telegram.SendMessage(ctx, text.PaidFeature[lang(ctx)], 0, p.settingsButton(ctx))
	if err != nil {
		log.Println("can't send premium feature message")
	}
}

func (p *Processor) processFile(ctx context.Context, message *tgbotapi.Message) {
	pieces := strings.Split(message.Document.FileName, ".")
	switch pieces[len(pieces)-1] {
	case "png", "jpg", "jpeg":
		message.Photo = []tgbotapi.PhotoSize{{FileID: message.Document.FileID}}
		p.process(ctx, message, "file")
		return
	}

	messageID, err := p.telegram.SendMessage(ctx, text.Loading[lang(ctx)], message.MessageID, nil)
	if err != nil {
		log.Println("can't send loading message")
		return
	}

	isTyping := p.telegram.SetTyping(ctx)
	defer isTyping.Store(false)

	var Text string
	errMsg := p.service.Load(ctx, message.Document)
	if supported, found := strings.CutPrefix(errMsg, "supported file formats: "); found {
		Text = fmt.Sprintf(text.UnsupportedFormat[lang(ctx)], supported)
	} else if errMsg != "" {
		log.Println("can't load file:", errMsg)
		Text = text.RequestFailed[lang(ctx)]
	} else {
		Text = fmt.Sprintf(text.FileLoaded[lang(ctx)], message.Document.FileName)
	}

	err = p.telegram.EditMessage(ctx, Text, messageID, nil)
	if err != nil {
		log.Println("can't edit process file message")
	}
}
