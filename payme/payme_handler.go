package payme

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/postgres"
	"github.com/gin-gonic/gin"
)

func Init() {

	merchantID, ok := os.LookupEnv("MERCHANT_ID")
	if !ok {
		log.Fatalf("merchant id is not specified")
	}
	constants.MerchantID = merchantID

	testKey, ok := os.LookupEnv("TEST_KEY")
	if !ok {
		log.Fatalf("test key is not specified")
	}
	constants.TestKey = testKey

	realKey, ok := os.LookupEnv("REAL_KEY")
	if !ok {
		log.Fatalf("real key is not specified")
	}
	constants.RealKey = realKey
}

func Handler(c *gin.Context) {

	if !functions.Authorized(c) {
		return
	}

	request := &types.PaymeRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		log.Printf("can't bind json: %v", err)
		c.JSON(200, gin.H{
			"error": gin.H{
				"code":    -32700,
				"message": "Parse error",
			},
		})
		return
	}

	var code int
	var response gin.H
	switch request.Method {
	case "CheckPerformTransaction":
		response, code = postgres.CheckPerformTransaction(&request.Params)
	case "CreateTransaction":
		response, code = postgres.CheckPerformTransaction(&request.Params)
		if code == 0 {
			response, code = postgres.CreateTransaction(&request.Params)
		}
	case "PerformTransaction":
		response, code = postgres.PerformTransaction(&request.Params)
	case "CancelTransaction":
		response, code = postgres.CancelTransaction(&request.Params)
	case "CheckTransaction":
		response, code = postgres.CheckTransaction(&request.Params)
	case "GetStatement":
		response, code = postgres.GetStatement(&request.Params)
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
		case -32400:
			message = "System error"
		}
		response = gin.H{
			"error": gin.H{
				"code":    code,
				"message": message,
			},
		}
	} else {
		response = gin.H{
			"result": response,
		}
	}

	response["id"] = request.ID
	c.JSON(200, response)
}

func CheckoutURL(ctx context.Context, amount int, Type string) string {

	userID := ctx.Value("user_id").(int64)

	orderID, err := postgres.NewOrder(userID, amount, Type)
	if err != nil {
		return ""
	}

	str := fmt.Sprintf("m=%s;ac.order_id=%d;a=%d", constants.MerchantID, orderID, amount)
	buffer := bytes.NewBuffer([]byte{})
	writer := base64.NewEncoder(base64.StdEncoding, buffer)
	_, err = writer.Write([]byte(str))
	if err != nil {
		return ""
	}
	_ = writer.Close()

	return "https://checkout.paycom.uz/" + buffer.String()
}
