package postgres

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/dro14/yordamchi/clients/telegram"
	"github.com/dro14/yordamchi/storage/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Postgres struct {
	telegram *telegram.Telegram
	redis    *redis.Redis
}

func New() *Postgres {
	url, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("database url is not specified")
	}

	if db == nil {
		var err error
		db, err = sql.Open("postgres", url)
		if err != nil {
			log.Fatal("can't connect to database:", err)
		}
	}

	return &Postgres{
		telegram: telegram.New(),
		redis:    redis.New(),
	}
}

func (p *Postgres) Stats() *gin.H {
	query := `SELECT
(SELECT COUNT(id) FROM users WHERE is_active = true),
(SELECT COUNT(DISTINCT user_id) FROM messages WHERE DATE(created_on) = $1),
(SELECT COUNT(id) FROM messages_legacy),
(SELECT COUNT(id) FROM messages),
(SELECT COUNT(id) FROM transactions WHERE state = 2);`

	var totalActiveUsers, dailyActiveUsers, legacyMessages, newMessages, totalPayments int
	args := []any{time.Now().AddDate(0, 0, -1).Format(time.DateOnly)}

	err := p.queryPayme(query, args, &totalActiveUsers, &dailyActiveUsers, &legacyMessages, &newMessages, &totalPayments)
	if err != nil {
		log.Printf("failed to get stats: %s", err)
		return nil
	}

	return &gin.H{
		"total_active_users": totalActiveUsers,
		"daily_active_users": dailyActiveUsers,
		"total_messages":     legacyMessages + newMessages,
		"total_payments":     totalPayments,
	}
}
