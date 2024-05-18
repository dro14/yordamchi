package postgres

import (
	"database/sql"
	"log"
	"os"

	"github.com/dro14/yordamchi/clients/telegram"
	"github.com/dro14/yordamchi/storage/redis"
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
