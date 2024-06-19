package click

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/dro14/yordamchi/payment/click/types"
	"github.com/gin-gonic/gin"
)

func (c *Click) CheckoutURL(ctx context.Context, amount int, Type string) string {
	orderID, err := c.postgres.NewOrder(id(ctx), amount, Type)
	if err != nil {
		return c.url
	}

	return fmt.Sprintf(c.url, c.serviceID, c.merchantID, amount/100, orderID)
}

func (c *Click) SingString(request *types.Request) string {
	var str string
	if request.MerchantPrepareID == 0 {
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

func (c *Click) Cancel(request *types.Request) gin.H {
	response := c.postgres.CheckOrder(request.MerchantTransID, request.Amount)
	if response != nil {
		return response
	}

	response = c.postgres.UpdateClickTransaction(request)
	if response != nil {
		return response
	}

	return gin.H{
		"error":      -9,
		"error_note": "Transaction cancelled",
	}
}

func (c *Click) Prepare(request *types.Request) gin.H {
	response := c.postgres.CheckOrder(request.MerchantTransID, request.Amount)
	if response != nil {
		return response
	}

	return c.postgres.CreateClickTransaction(request)
}

func (c *Click) Complete(request *types.Request) gin.H {
	response := c.postgres.CheckOrder(request.MerchantTransID, request.Amount)
	if response != nil {
		return response
	}

	response = c.postgres.UpdateClickTransaction(request)
	if response != nil {
		return response
	}

	return gin.H{
		"click_trans_id":      request.ClickTransID,
		"merchant_trans_id":   request.MerchantTransID,
		"merchant_confirm_id": request.MerchantPrepareID,
		"error":               0,
		"error_note":          "Success",
	}
}
