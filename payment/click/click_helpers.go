package click

import (
	"crypto/md5"
	"fmt"
	"github.com/dro14/yordamchi/payment/click/types"
)

func (c *Click) singString(request *types.Request, isPrepare bool) string {
	var str string
	if isPrepare {
		str = fmt.Sprintf("%d%d%s%d%.0f%d%s",
			request.ClickTransID,
			request.ServiceID,
			c.secretKey,
			request.MerchantTransID,
			request.Amount,
			request.Action,
			request.SignTime,
		)
	} else {
		str = fmt.Sprintf("%d%d%s%d%d%.0f%d%s",
			request.ClickTransID,
			request.ServiceID,
			c.secretKey,
			request.MerchantTransID,
			request.MerchantPrepareID,
			request.Amount,
			request.Action,
			request.SignTime,
		)
	}
	hash := md5.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
