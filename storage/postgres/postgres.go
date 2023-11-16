package postgres

import (
	"database/sql"
	"github.com/dro14/yordamchi/clients/telegram"
	"github.com/dro14/yordamchi/storage/redis"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db       *sql.DB
	telegram *telegram.Telegram
	redis    *redis.Redis
}

func New() *Postgres {
	url, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("database url is not specified")
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("can't connect to database:", err)
	}

	return &Postgres{
		db:       db,
		telegram: telegram.New(),
		redis:    redis.New(),
	}
}
