package click

import (
	"github.com/dro14/yordamchi/payment/click/types"
	"github.com/gin-gonic/gin"
)

func (c *Click) Prepare(request *types.Request) gin.H {
	if request.SignString != c.singString(request, true) {
		return gin.H{"error": -1, "error_note": "SIGN CHECK FAILED!"}
	} else if h := c.postgres.CheckOrder(request.MerchantTransID, request.Amount); h != nil {
		return h
	}

	return c.postgres.CreateClickTransaction(request)
}

func (c *Click) Complete(request *types.Request) gin.H {
	if request.SignString != c.singString(request, false) {
		return gin.H{"error": -1, "error_note": "SIGN CHECK FAILED!"}
	} else if h := c.postgres.CheckOrder(request.MerchantTransID, request.Amount); h != nil {
		return h
	}

	response := c.postgres.UpdateClickTransaction(request, true)
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

func (c *Click) Cancel(request *types.Request) gin.H {
	if request.SignString != c.singString(request, false) {
		return gin.H{"error": -1, "error_note": "SIGN CHECK FAILED!"}
	}

	response := c.postgres.UpdateClickTransaction(request, false)
	if response != nil {
		return response
	}

	return gin.H{
		"error":      -9,
		"error_note": "Transaction cancelled",
	}
}
