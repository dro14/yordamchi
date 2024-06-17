package postgres

import (
	"log"
	"time"
)

func (p *Postgres) NewOrder(userID int64, amount int, Type string) (int, error) {
	query := "SELECT id, created_at, updated_at FROM orders WHERE user_id = $1 AND amount = $2 AND type = $3 ORDER BY created_at DESC;"
	args := []any{userID, amount, Type}
	var id int
	var createdAt, updatedAt string
	err := p.queryPayment(query, args, &id, &createdAt, &updatedAt)
	if err != nil || createdAt != updatedAt {
		query = "INSERT INTO orders (user_id, amount, type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
		args = []any{userID, amount, Type, time.Now().Format(time.DateTime), time.Now().Format(time.DateTime)}
		err = p.queryPayment(query, args, &id)
		if err != nil {
			log.Println("can't create order:", err)
			return 0, err
		}
	}
	return id, nil
}
