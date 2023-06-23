package telegram

import (
	"context"
	"fmt"

	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
)

func settingsMessage(ctx context.Context, lang string) string {

	balance, userStatus := settings(ctx)

	switch userStatus {
	case types.PremiumStatus:
		if lang == "uz" {
			return fmt.Sprintf(text.Settings[lang], "", balance, "premium")
		} else if lang == "ru" {
			return fmt.Sprintf(text.Settings[lang], "", balance, "премиум")
		} else {
			return fmt.Sprintf(text.Settings[lang], "", balance, "premium")
		}
	case types.FreeStatus:
		if lang == "uz" {
			return fmt.Sprintf(text.Settings[lang], "Bugunga, ", balance, "bepul")
		} else if lang == "ru" {
			return fmt.Sprintf(text.Settings[lang], "На сегодня, ", balance, "бесплатных")
		} else {
			return fmt.Sprintf(text.Settings[lang], "For today, ", balance, "free")
		}
	case types.ExhaustedStatus:
		return text.Exhausted[lang]
	}

	return ""
}

func settings(ctx context.Context) (int, types.UserStatus) {

	userStatus, _ := redis.Status(ctx)

	switch userStatus {
	case types.PremiumStatus, types.FreeStatus, types.ExhaustedStatus:
		requests, err := redis.Balance(ctx)
		if err != nil {
			return -1, types.UnknownStatus
		}
		return requests, userStatus
	default:
		return -1, types.UnknownStatus
	}
}
