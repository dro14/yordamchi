package postgres

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dro14/yordamchi/payme/types"
	"github.com/dro14/yordamchi/processor/text"
	"github.com/gin-gonic/gin"
)

func (p *Postgres) NewOrder(userID int64, amount int, Type string) (int, error) {
	query := "SELECT id, updated_at FROM orders WHERE user_id = $1 AND amount = $2 AND type = $3 ORDER BY created_at DESC;"
	args := []any{userID, amount, Type}
	var id int
	var updatedAt string
	err := p.queryPayme(query, args, &id, &updatedAt)
	if err != nil || updatedAt != "" {
		query = "INSERT INTO orders (user_id, amount, type, created_at) VALUES ($1, $2, $3, $4) RETURNING id;"
		args = []any{userID, amount, Type, time.Now().Format(time.DateTime)}
		err = p.queryPayme(query, args, &id)
		if err != nil {
			log.Println("can't create order:", err)
			return 0, err
		}
	}
	return id, nil
}

func (p *Postgres) CheckPerformTransaction(params *types.Params) (gin.H, int) {
	query := "SELECT amount, type, user_id FROM orders WHERE id = $1;"
	args := []any{params.Account.OrderID}
	var amount int
	var order string
	var userID int64
	err := p.queryPayme(query, args, &amount, &order, &userID)
	if err != nil {
		log.Println("can't get order:", err)
		return nil, -31050
	}

	if amount != params.Amount {
		log.Printf("invalid amount: %d != %d", amount, params.Amount)
		return nil, -31001
	}

	subscription, orderType, found := strings.Cut(order, ":")
	if !found {
		log.Println("invalid order:", order)
		return nil, -31052
	}

	var title string
	switch orderType {
	case "premium", "gpt-4":
		switch subscription {
		case "daily":
			title = "Дневная премиум подписка"
		case "weekly":
			title = "Недельная премиум подписка"
		case "monthly":
			title = "Месячная премиум подписка"
		default:
			log.Printf("user %d: invalid subscripion: %v", userID, order)
			return nil, -31052
		}
	case "unlimited":
		switch subscription {
		case "weekly":
			title = "Недельная безлимитная подписка"
		case "monthly":
			title = "Месячная безлимитная подписка"
		default:
			log.Printf("user %d: invalid subscripion: %v", userID, order)
			return nil, -31052
		}
	case "images", "dall-e-3":
		_, err = strconv.Atoi(subscription)
		if err != nil {
			log.Println("invalid number of images:", subscription)
			return nil, -31052
		}
		title = fmt.Sprintf("%s изображений DALL-E 3", subscription)
	default:
		log.Println("invalid order type:", orderType)
		return nil, -31052
	}

	result := gin.H{
		"allow": true,
		"detail": gin.H{
			"receipt_type": 0,
			"items": []gin.H{
				{
					"title":        title,
					"price":        amount,
					"count":        1,
					"code":         "10318001001000000",
					"package_code": "1501319",
					"vat_percent":  0,
				},
			},
		},
	}
	return result, 0
}

func (p *Postgres) CreateTransaction(params *types.Params) (gin.H, int) {
	result, code := p.CheckTransaction(params)
	if code != -31003 {
		if code == 0 {
			if result["state"] == 1 {
				delete(result, "perform_time")
				delete(result, "cancel_time")
				delete(result, "reason")
				return result, 0
			} else {
				log.Println("invalid state:", result["state"])
				return nil, -31008
			}
		} else {
			return nil, code
		}
	}

	result = gin.H{
		"create_time": time.Now().UnixMilli(),
		"transaction": transactionHash(params),
		"state":       1,
	}

	query := "INSERT INTO transactions (id, time, amount, order_id, create_time, transaction, state) VALUES ($1, $2, $3, $4, $5, $6, $7);"
	args := []any{params.ID, params.Time, params.Amount, params.Account.OrderID, result["create_time"], result["transaction"], result["state"]}
	err := p.execPayme(query, args)
	if err != nil {
		log.Println("can't create transaction:", err)
		return nil, -32400
	}
	return result, 0
}

func (p *Postgres) PerformTransaction(params *types.Params) (gin.H, int) {
	result, code := p.CheckTransaction(params)
	if code != 0 {
		return nil, code
	}

	if result["state"] == 1 {
		result["state"] = 2
		result["perform_time"] = time.Now().UnixMilli()
	} else if result["state"] == 2 {
		log.Printf("transaction %s already performed", params.ID)
	} else {
		log.Printf("invalid state of transaction %s: %s", params.ID, result["state"])
		return nil, -31008
	}

	query := "UPDATE transactions SET state = $1, perform_time = $2 WHERE id = $3 RETURNING order_id;"
	args := []any{result["state"], result["perform_time"], params.ID}
	var orderID int
	err := p.queryPayme(query, args, &orderID)
	if err != nil {
		log.Printf("can't perform transaction %s: %s", params.ID, err)
		return nil, -32400
	}

	query = "UPDATE orders SET updated_at = $1 WHERE id = $2 RETURNING user_id, type;"
	args = []any{time.Now().Format(time.DateTime), orderID}
	var userID int64
	var Type string
	err = p.queryPayme(query, args, &userID, &Type)
	if err != nil {
		log.Printf("can't update order %d: %s", orderID, err)
		return nil, -32400
	}

	ctx := context.WithValue(context.Background(), "user_id", userID)
	err = p.redis.PerformTransaction(ctx, Type)
	if err != nil {
		log.Printf("user %d: can't perform transaction: %s", userID, err)
		return nil, -32400
	}

	ctx, _ = p.redis.Lang(ctx, "uz")
	_, err = p.telegram.SendMessage(ctx, text.Success[lang(ctx)], 0, nil)
	if err != nil {
		log.Printf("user %d: can't send success message: %s", userID, err)
	}

	delete(result, "create_time")
	delete(result, "cancel_time")
	delete(result, "reason")
	return result, 0
}

func (p *Postgres) CancelTransaction(params *types.Params) (gin.H, int) {
	result, code := p.CheckTransaction(params)
	if code != 0 {
		return nil, code
	}

	if result["state"] == 1 {
		result["state"] = -1
		result["cancel_time"] = time.Now().UnixMilli()
	} else if result["state"] == 2 {
		result["state"] = -2
		result["cancel_time"] = time.Now().UnixMilli()
	} else if result["state"] == -1 || result["state"] == -2 {
		delete(result, "create_time")
		delete(result, "perform_time")
		delete(result, "reason")
		return result, 0
	} else {
		log.Println("invalid state:", result["state"])
		return nil, -31008
	}

	query := "UPDATE transactions SET state = $1, cancel_time = $2, reason = $3 WHERE id = $4 RETURNING order_id;"
	args := []any{result["state"], result["cancel_time"], params.Reason, params.ID}
	var orderID int
	err := p.queryPayme(query, args, &orderID)
	if err != nil {
		log.Println("can't cancel transaction:", err)
		return nil, -32400
	}

	query = "UPDATE orders SET updated_at = $1 WHERE id = $2;"
	args = []any{time.Now().Format(time.DateTime), orderID}
	err = p.execPayme(query, args)
	if err != nil {
		log.Println("can't update order:", err)
		return nil, -32400
	}

	delete(result, "create_time")
	delete(result, "perform_time")
	delete(result, "reason")
	return result, 0
}

func (p *Postgres) CheckTransaction(params *types.Params) (gin.H, int) {
	if params.Account.OrderID == "" {
		params.Account.OrderID = "0"
	}

	query := "SELECT id, create_time, transaction, state, perform_time, cancel_time, reason FROM transactions WHERE id = $1 OR order_id = $2;"
	args := []any{params.ID, params.Account.OrderID}
	var id string
	var createTime int64
	var transaction string
	var state int
	var performTime int64
	var cancelTime int64
	var reason int
	err := p.queryPayme(query, args, &id, &createTime, &transaction, &state, &performTime, &cancelTime, &reason)
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

func (p *Postgres) GetStatement(params *types.Params) (gin.H, int) {
	rows, err := db.Query(`SELECT * FROM transactions WHERE time >= $1 AND time <= $2 ORDER BY time;`, params.From, params.To)
	if err != nil {
		log.Println("can't get transactions:", err)
		return nil, -32400
	}

	var id string
	var Time int64
	var amount int
	var orderID string
	var createTime int64
	var transaction string
	var state int
	var performTime int64
	var cancelTime int64
	var reason int
	var transactions []gin.H

	for rows.Next() {
		err = rows.Scan(&id, &Time, &amount, &orderID, &createTime, &transaction, &state, &performTime, &cancelTime, &reason)
		if err != nil {
			log.Println("can't get transaction:", err)
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

	if len(transactions) == 0 {
		return gin.H{"transactions": []gin.H{}}, 0
	}
	return gin.H{"transactions": transactions}, 0
}
