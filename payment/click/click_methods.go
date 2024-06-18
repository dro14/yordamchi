package click

import (
	"github.com/dro14/yordamchi/payment/click/types"
	"github.com/gin-gonic/gin"
)

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
