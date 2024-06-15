package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/dro14/yordamchi/payment/payme/types"
)

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

func transactionHash(params *types.Params) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s_%d_%d", params.Account.OrderID, params.Time, params.Amount)))
	hashValue := hash.Sum(nil)
	return hex.EncodeToString(hashValue)
}
