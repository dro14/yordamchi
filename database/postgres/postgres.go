package postgres

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/dro14/yordamchi/lib/types"
	"github.com/gotd/td/tg"
	_ "github.com/lib/pq"
)

type Database struct {
	Postgres *sql.DB
}

func New() *Database {

	url, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatalf("database url is not specified")
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}

	return &Database{
		Postgres: db,
	}
}

func (db *Database) JoinUser(ctx context.Context, user *tg.User, joinedBy int64) {

	var (
		attempts   int
		retryDelay = constants.RetryDelay
		lang       = ctx.Value("language_code").(string)
		date       = int64(ctx.Value("date").(int))
		joinedAt   = time.Unix(date, 0).Format("2006-01-02 15:04:05")
	)

Retry:
	attempts++
	_, _ = db.Postgres.ExecContext(ctx, `
		INSERT INTO users
    		(id, first_name, last_name, username, language_code)
		VALUES
			($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING;`,
		user.ID, user.FirstName, user.LastName, user.Username, lang)

	_, err := db.Postgres.ExecContext(ctx, `
		INSERT INTO user_configs
			(id, joined_by, is_active, joined_at)
		VALUES
		    ($1, $2, $3, $4);`,
		user.ID, joinedBy, true, joinedAt)

	if err != nil {
		if strings.Contains(err.Error(), e.UniqueConstraint) {
			if !db.IsActive(ctx, user) {
				db.RejoinUser(ctx, user)
			}
		} else {
			log.Printf("can't join user: %v", err)
			if attempts < constants.RetryAttempts {
				functions.Sleep(&retryDelay)
				goto Retry
			} else {
				log.Printf("joining user failed after %d attempts", attempts)
			}
		}
	}
}

func (db *Database) SaveMessage(ctx context.Context, stats *types.Stats, user *tg.User) {

	var (
		attempts    int
		retryDelay  = constants.RetryDelay
		lang        = ctx.Value("language_code").(string)
		date        = int64(ctx.Value("date").(int))
		createdOn   = time.Unix(date, 0).Format("2006-01-02")
		promptedAt  = time.Unix(date, 0).Format("15:04:05")
		completedAt = time.Unix(stats.CompletedAt, 0).Format("15:04:05")
	)

	if !db.IsActive(ctx, user) {
		db.RejoinUser(ctx, user)
	}

Retry:
	attempts++
	_, err := db.Postgres.ExecContext(ctx, `
		INSERT INTO messages
			(user_id, is_premium, created_on, prompted_at, completed_at, first_send, last_edit, prompt_tokens, prompt_length, completion_tokens, completion_length, activity, requests, attempts, finish_reason, language_code)
		VAlUES
		    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);`,
		user.ID, stats.IsPremium, createdOn, promptedAt, completedAt, stats.FirstSend, stats.LastEdit, stats.PromptTokens, stats.PromptLength, stats.CompletionTokens, stats.CompletionLength, stats.Activity, stats.Requests, stats.Attempts, stats.FinishReason, lang)

	if err != nil {
		if strings.Contains(err.Error(), e.ForeignKeyConstraint) {
			db.JoinUser(ctx, user, 0)
			goto Retry
		} else {
			log.Printf("can't save message: %v", err)
			if attempts < constants.RetryAttempts {
				functions.Sleep(&retryDelay)
				goto Retry
			} else {
				log.Printf("saving message failed after %d attempts", attempts)
			}
		}
	}
}

func (db *Database) DeactivateUser(ctx context.Context, user *tg.User) {

	var (
		attempts      int
		retryDelay    = constants.RetryDelay
		date          = int64(ctx.Value("date").(int))
		deactivatedAt = time.Unix(date, 0).Format("2006-01-02 15:04:05")
	)

Retry:
	attempts++
	_, err := db.Postgres.ExecContext(ctx, `
	UPDATE user_configs
	SET
		is_active = $1,
	    deactivated_at = $2
	WHERE
	    id = $3;`,
		false, deactivatedAt, user.ID)

	if err != nil {
		log.Printf("can't deactivate user: %v", err)
		if attempts < constants.RetryAttempts {
			functions.Sleep(&retryDelay)
			goto Retry
		} else {
			log.Printf("deactivating user failed after %d attempts", attempts)
		}
	}
}

func (db *Database) RejoinUser(ctx context.Context, user *tg.User) {

	var (
		attempts   int
		retryDelay = constants.RetryDelay
		date       = int64(ctx.Value("date").(int))
		rejoinedAt = time.Unix(date, 0).Format("2006-01-02 15:04:05")
	)

Retry:
	attempts++
	_, err := db.Postgres.ExecContext(ctx, `
	UPDATE user_configs
	SET
	    is_active = $1,
	    rejoined_at = $2
	WHERE
	    id = $3;`,
		true, rejoinedAt, user.ID)

	if err != nil {
		log.Printf("can't rejoin user: %v", err)
		if attempts < constants.RetryAttempts {
			functions.Sleep(&retryDelay)
			goto Retry
		} else {
			log.Printf("rejoining user failed after %d attempts", attempts)
		}
	}
}
