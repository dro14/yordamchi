package telegram

import (
	"context"
	"log"

	"github.com/dro14/yordamchi/client/telegram"
	"github.com/dro14/yordamchi/processor/telegram/button"
	"github.com/dro14/yordamchi/text"
)

func exhausted(ctx context.Context) {

	_, err := telegram.SendMessage(ctx, text.Exhausted[lang(ctx)], 0, button.Premium(ctx, lang(ctx)))
	if err != nil {
		log.Printf("can't send exhausted message")
	}
}
