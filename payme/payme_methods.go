package payme

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/dro14/yordamchi/payme/types"
	"github.com/gin-gonic/gin"
)

func (p *Payme) Respond(c *gin.Context, request *types.Request) gin.H {
	response := p.authorized(c)
	if response == nil {
		var code int
		switch request.Method {
		case "CheckPerformTransaction":
			response, code = p.postgres.CheckPerformTransaction(&request.Params)
		case "CreateTransaction":
			response, code = p.postgres.CheckPerformTransaction(&request.Params)
			if code == 0 {
				response, code = p.postgres.CreateTransaction(&request.Params)
			}
		case "PerformTransaction":
			response, code = p.postgres.PerformTransaction(&request.Params)
		case "CancelTransaction":
			response, code = p.postgres.CancelTransaction(&request.Params)
		case "CheckTransaction":
			response, code = p.postgres.CheckTransaction(&request.Params)
		case "GetStatement":
			response, code = p.postgres.GetStatement(&request.Params)
		}

		if code != 0 {
			var message string
			switch code {
			case -31001:
				message = "Invalid amount"
			case -31003:
				message = "Transaction not found"
			case -31008:
				message = "Impossible to complete operation"
			case -31050:
				message = "Invalid account: order_id"
			case -31051:
				message = "Invalid id"
			case -31052:
				message = "Invalid account: type"
			case -32400:
				message = "System error"
			}
			response = gin.H{"error": gin.H{"code": code, "message": message}}
		} else {
			response = gin.H{"result": response}
		}
	}
	response["id"] = request.ID
	return response
}

func (p *Payme) CheckoutURL(ctx context.Context, amount int, Type string) string {
	URL := "https://checkout.paycom.uz/"
	userID := ctx.Value("user_id").(int64)
	orderID, err := p.postgres.NewOrder(userID, amount, Type)
	if err != nil {
		return URL
	}
	buffer := bytes.NewBuffer([]byte{})
	writer := base64.NewEncoder(base64.StdEncoding, buffer)
	s := fmt.Sprintf("m=%s;ac.order_id=%d;a=%d;l=%s", p.merchantID, orderID, amount, lang(ctx))
	_, err = writer.Write([]byte(s))
	if err != nil {
		return URL
	}
	_ = writer.Close()
	return URL + buffer.String()
}
