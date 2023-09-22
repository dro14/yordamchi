package postgres

import (
	"context"
	"log"
	"time"

	"github.com/dro14/yordamchi/client/telegram"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/dro14/yordamchi/redis"
	"github.com/dro14/yordamchi/text"
	"github.com/gin-gonic/gin"
)

func NewOrder(userID int64, amount int, Type string) (int, error) {

	var id int
	err := db.QueryRow(
		`INSERT INTO orders_test
    	(user_id, amount, type, created_at)
		VALUES
		($1, $2, $3, $4) RETURNING id;`,
		userID, amount, Type, time.Now().Format(time.DateTime)).
		Scan(&id)
	if err != nil {
		log.Printf("can't create order: %v", err)
		return 0, err
	}

	return id, nil
}

func CheckPerformTransaction(params *types.Params) (gin.H, int) {

	var (
		amount int
		Type   string
	)

	err := db.QueryRow(`SELECT amount, type FROM orders_test WHERE id = $1;`, params.Account.OrderID).Scan(&amount, &Type)
	if err != nil {
		log.Printf("can't get order: %v", err)
		return nil, -31050
	}

	if amount != params.Amount {
		log.Printf("invalid amount: %d != %d", amount, params.Amount)
		return nil, -31001
	}

	var (
		title string
		price int
		count int
	)

	switch Type {
	case "weekly":
		title = "Weekly tariff for ChatGPT"
		price = amount
		count = 1
	case "monthly":
		title = "Monthly tariff for ChatGPT"
		price = amount
		count = 1
	case "gpt-4":
		title = "Tokens for GPT-4"
		price = 100
		count = amount / 100
	}

	result := gin.H{
		"allow": true,
		"detail": gin.H{
			"receipt_type": 0,
			"items": []gin.H{
				{
					"title":        title,
					"price":        price,
					"count":        count,
					"code":         10318001001000000,
					"package_code": 1501319,
					"vat_percent":  0,
				},
			},
		},
	}

	return result, 0
}

func CreateTransaction(params *types.Params) (gin.H, int) {

	result, code := CheckTransaction(params)
	if code != -31003 {
		if code == 0 {
			if result["state"] == 1 {
				delete(result, "perform_time")
				delete(result, "cancel_time")
				delete(result, "reason")
				return result, 0
			} else {
				log.Printf("invalid state: %d", result["state"])
				return nil, -31008
			}
		} else {
			return nil, code
		}
	}

	result = gin.H{
		"create_time": time.Now().UnixMilli(),
		"transaction": functions.Transaction(params),
		"state":       1,
	}

	_, err := db.Exec(
		`INSERT INTO transactions_test
    	(id, time, amount, order_id, create_time, transaction, state)
		VALUES
    	($1, $2, $3, $4, $5, $6, $7);`,
		params.ID, params.Time, params.Amount, params.Account.OrderID, result["create_time"], result["transaction"], result["state"])
	if err != nil {
		log.Printf("can't create transaction: %v", err)
		return nil, -32400
	}

	return result, 0
}

func PerformTransaction(params *types.Params) (gin.H, int) {

	result, code := CheckTransaction(params)
	if code != 0 {
		return nil, code
	}
	delete(result, "create_time")
	delete(result, "cancel_time")
	delete(result, "reason")

	if result["state"] == 1 {
		result["state"] = 2
		result["perform_time"] = time.Now().UnixMilli()
	} else if result["state"] == 2 {
		log.Printf("transaction already performed: %s", params.ID)
	} else {
		log.Printf("invalid state: %d", result["state"])
		return nil, -31008
	}

	var orderID int
	err := db.QueryRow(`UPDATE transactions_test SET state = $1, perform_time = $2 WHERE id = $3 RETURNING order_id;`, result["state"], result["perform_time"], params.ID).Scan(&orderID)
	if err != nil {
		log.Printf("can't perform transaction: %v", err)
		return nil, -32400
	}

	var (
		userID int64
		amount int
		Type   string
	)

	err = db.QueryRow(`UPDATE orders_test SET updated_at = $1 WHERE id = $2 RETURNING user_id, amount, type;`, time.Now().Format("2006-01-02 15:04:05"), orderID).Scan(&userID, &amount, &Type)
	if err != nil {
		log.Printf("can't update order: %v", err)
		return nil, -32400
	}

	err = redis.PerformTransaction(userID, amount, Type)
	if err != nil {
		return nil, -32400
	}

	var lang string
	err = db.QueryRow(`SELECT language_code FROM users WHERE id = $1;`, userID).Scan(&lang)
	if err != nil {
		log.Printf("can't get language_code: %v", err)
		lang = "uz"
	}

	ctx := context.WithValue(context.Background(), "user_id", userID)
	_, err = telegram.SendMessage(ctx, text.Success[lang], 0, nil)
	if err != nil {
		log.Printf("can't send success message: %v", err)
	}

	return result, 0
}

func CancelTransaction(params *types.Params) (gin.H, int) {

	result, code := CheckTransaction(params)
	if code != 0 {
		return nil, code
	}
	delete(result, "create_time")
	delete(result, "perform_time")
	delete(result, "reason")

	if result["state"] == 1 {
		result["state"] = -1
		result["cancel_time"] = time.Now().UnixMilli()
	} else if result["state"] == 2 {
		result["state"] = -2
		result["cancel_time"] = time.Now().UnixMilli()
	} else if result["state"] == -1 || result["state"] == -2 {
		return result, 0
	} else {
		log.Printf("invalid state: %d", result["state"])
		return nil, -31008
	}

	var orderID int
	err := db.QueryRow(`UPDATE transactions_test SET state = $1, cancel_time = $2, reason = $3 WHERE id = $4 RETURNING order_id;`, result["state"], result["cancel_time"], params.Reason, params.ID).Scan(&orderID)
	if err != nil {
		log.Printf("can't cancel transaction: %v", err)
		return nil, -32400
	}

	_, err = db.Exec(`UPDATE orders_test SET updated_at = $1 WHERE id = $2;`, time.Now().Format("2006-01-02 15:04:05"), orderID)
	if err != nil {
		log.Printf("can't update order: %v", err)
		return nil, -32400
	}

	return result, 0
}

func CheckTransaction(params *types.Params) (gin.H, int) {

	var (
		id          string
		createTime  int64
		transaction string
		state       int
		performTime int64
		cancelTime  int64
		reason      int
	)

	if params.Account.OrderID == "" {
		params.Account.OrderID = "0"
	}

	err := db.QueryRow(
		`SELECT id, create_time, transaction, state, perform_time, cancel_time, reason
		FROM transactions_test
		WHERE id = $1 OR order_id = $2;`, params.ID, params.Account.OrderID).
		Scan(&id, &createTime, &transaction, &state, &performTime, &cancelTime, &reason)
	if err != nil {
		return nil, -31003
	}

	if id != params.ID {
		log.Printf("invalid id: %s != %s", id, params.ID)
		return nil, -31051
	}

	result := gin.H{
		"create_time":  createTime,
		"transaction":  transaction,
		"state":        state,
		"perform_time": performTime,
		"cancel_time":  cancelTime,
		"reason":       reason,
	}

	if reason == 0 {
		result["reason"] = nil
	}

	return result, 0
}

func GetStatement(params *types.Params) (gin.H, int) {

	rows, err := db.Query(`SELECT * FROM transactions_test WHERE time >= $1 AND time <= $2 ORDER BY time;`, params.From, params.To)
	if err != nil {
		log.Printf("can't get transactions: %v", err)
		return nil, -32400
	}

	var (
		id          string
		Time        int64
		amount      int
		orderID     string
		createTime  int64
		transaction string
		state       int
		performTime int64
		cancelTime  int64
		reason      int
	)

	var transactions []gin.H

	for rows.Next() {

		err = rows.Scan(&id, &Time, &amount, &orderID, &createTime, &transaction, &state, &performTime, &cancelTime, &reason)
		if err != nil {
			log.Printf("can't get transaction: %v", err)
			return nil, -32400
		}

		result := gin.H{
			"id":     id,
			"time":   Time,
			"amount": amount,
			"account": gin.H{
				"order_id": orderID,
			},
			"create_time":  createTime,
			"transaction":  transaction,
			"state":        state,
			"perform_time": performTime,
			"cancel_time":  cancelTime,
			"reason":       reason,
		}

		if reason == 0 {
			result["reason"] = nil
		}

		transactions = append(transactions, result)
	}

	if transactions != nil {
		return gin.H{
			"transactions": transactions,
		}, 0
	} else {
		return gin.H{
			"transactions": []gin.H{},
		}, 0
	}
}
