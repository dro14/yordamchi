package postgres

import (
	"context"
	"log"
	"time"

	"github.com/dro14/yordamchi/payment/click/methods"
	"github.com/dro14/yordamchi/payment/click/types"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/gin-gonic/gin"
)

func (p *Postgres) CheckOrder(orderID int, orderAmount float64) gin.H {
	query := "SELECT amount FROM orders WHERE id = $1;"
	args := []any{orderID}
	var amount int
	err := p.queryPayment(query, args, &amount)
	if err != nil {
		log.Println("can't get order:", err)
		return gin.H{"error": -5, "error_note": "User does not exist"}
	}
	if amount != int(orderAmount*100) {
		log.Printf("invalid amount: %d != %d", amount, int(orderAmount*100))
		return gin.H{"error": -2, "error_note": "Incorrect parameter amount"}
	}
	return nil
}

func (p *Postgres) CreateClickTransaction(request *types.Request) gin.H {
	query := "INSERT INTO click_transactions (click_trans_id, service_id, click_paydoc_id, merchant_trans_id, amount, action, error, error_note, sign_time, sign_string) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id;"
	args := []any{request.ClickTransID, request.ServiceID, request.ClickPaydocID, request.MerchantTransID, request.Amount, request.Action, request.Error, request.ErrorNote, request.SignTime, request.SignString}
	var id int
	err := p.queryPayment(query, args, &id)
	if err != nil {
		log.Println("can't create click transaction:", err)
		return gin.H{"error": -5, "error_note": "User does not exist"}
	}

	return gin.H{
		"click_trans_id":      request.ClickTransID,
		"merchant_trans_id":   request.MerchantTransID,
		"merchant_prepare_id": id,
		"error":               0,
		"error_note":          "Success",
	}
}

func (p *Postgres) UpdateClickTransaction(request *types.Request) gin.H {
	query := "SELECT action FROM click_transactions WHERE id = $1;"
	args := []any{request.MerchantPrepareID}
	var action int
	err := p.queryPayment(query, args, &action)
	if err != nil {
		log.Println("can't get click transaction:", err)
		return gin.H{"error": -6, "error_note": "Transaction does not exist"}
	}

	if action == methods.Complete {
		log.Println("transaction has been performed:", request.MerchantPrepareID)
		return gin.H{"error": -4, "error_note": "Already paid"}
	} else if action == methods.Cancel {
		log.Println("transaction has been canceled:", request.MerchantPrepareID)
		return gin.H{"error": -9, "error_note": "Transaction cancelled"}
	}

	if request.Error == 0 {
		action = methods.Complete
	} else {
		action = methods.Cancel
	}

	query = "UPDATE click_transactions SET click_trans_id = $1, service_id = $2, click_paydoc_id = $3, merchant_trans_id = $4, amount = $5, action = $6, error = $7, error_note = $8, sign_time = $9, sign_string = $10 WHERE id = $11;"
	args = []any{request.ClickTransID, request.ServiceID, request.ClickPaydocID, request.MerchantTransID, request.Amount, action, request.Error, request.ErrorNote, request.SignTime, request.SignString, request.MerchantPrepareID}
	err = p.execPayment(query, args)
	if err != nil {
		log.Println("can't update click transaction:", err)
		return gin.H{"error": -7, "error_note": "Failed to update user"}
	}

	query = "UPDATE orders SET updated_at = $1 WHERE id = $2 RETURNING user_id, type;"
	args = []any{time.Now().Format(time.DateTime), request.MerchantTransID}
	var userID int64
	var Type string
	err = p.queryPayment(query, args, &userID, &Type)
	if err != nil {
		log.Printf("can't update order %d: %s", request.MerchantTransID, err)
		return gin.H{"error": -7, "error_note": "Failed to update user"}
	}

	ctx := context.WithValue(context.Background(), "user_id", userID)
	var message string
	if request.Error == 0 {
		err = p.redis.PerformTransaction(ctx, Type)
		if err != nil {
			log.Printf("user %d: can't perform transaction: %s", userID, err)
			return gin.H{"error": -7, "error_note": "Failed to update user"}
		}
		message = text.Success[lang(ctx)]
	} else {
		err = p.redis.CancelTransaction(ctx, Type)
		if err != nil {
			log.Printf("user %d: can't cancel transaction: %s", userID, err)
			return gin.H{"error": -7, "error_note": "Failed to update user"}
		}
		message = text.Cancel[lang(ctx)]
	}

	ctx, _ = p.redis.Lang(ctx, "uz")
	_, err = p.telegram.SendMessage(ctx, message, 0, nil)
	if err != nil {
		log.Printf("user %d: can't send message: %s", userID, err)
	}

	return nil
}
