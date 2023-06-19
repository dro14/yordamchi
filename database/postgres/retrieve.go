package postgres

import (
	"context"
	"log"
	"strings"

	"github.com/dro14/yordamchi/lib/constants"
	"github.com/dro14/yordamchi/lib/e"
	"github.com/dro14/yordamchi/lib/functions"
	"github.com/gotd/td/tg"
)

func (db *Database) IsActive(ctx context.Context, user *tg.User) bool {

	var (
		attempts   int
		retryDelay = constants.RetryDelay
		isActive   bool
	)

Retry:
	attempts++
	err := db.Postgres.QueryRowContext(ctx, `
		SELECT is_active
		FROM user_configs
		WHERE id = $1;`,
		user.ID).Scan(&isActive)

	if err != nil {
		if strings.Contains(err.Error(), e.NotFound) {
			db.JoinUser(ctx, user, 0)
			goto Retry
		} else {
			log.Printf("can't get is_active: %v", err)
			if attempts < constants.RetryAttempts {
				functions.Sleep(&retryDelay)
				goto Retry
			} else {
				log.Printf("getting is_active failed after %d attempts", attempts)
			}
		}
	}

	return isActive
}

func (db *Database) JoinedAt(ctx context.Context, user *tg.User) string {

	var (
		attempts   int
		retryDelay = constants.RetryDelay
		joinedAt   string
	)

Retry:
	attempts++
	err := db.Postgres.QueryRowContext(ctx, `
		SELECT joined_at
		FROM user_configs
		WHERE id = $1;`,
		user.ID).Scan(&joinedAt)

	if err != nil {
		if strings.Contains(err.Error(), e.NotFound) {
			db.JoinUser(ctx, user, 0)
			goto Retry
		} else {
			log.Printf("can't get joined_at: %v", err)
			if attempts < constants.RetryAttempts {
				functions.Sleep(&retryDelay)
				goto Retry
			} else {
				log.Printf("getting joined_at failed after %d attempts", attempts)
			}
		}
	}

	return joinedAt
}

func (db *Database) DeactivatedAt(ctx context.Context, user *tg.User) string {

	var (
		attempts      int
		retryDelay    = constants.RetryDelay
		deactivatedAt string
	)

Retry:
	attempts++
	err := db.Postgres.QueryRowContext(ctx, `
		SELECT deactivated_at
		FROM user_configs
		WHERE id = $1;`,
		user.ID).Scan(&deactivatedAt)

	if err != nil {
		if strings.Contains(err.Error(), e.UnsupportedConversion) {
			return ""
		} else if strings.Contains(err.Error(), e.NotFound) {
			db.JoinUser(ctx, user, 0)
			goto Retry
		} else {
			log.Printf("can't get deactivated_at: %v", err)
			if attempts < constants.RetryAttempts {
				functions.Sleep(&retryDelay)
				goto Retry
			} else {
				log.Printf("getting deactivated_at failed after %d attempts", attempts)
			}
		}
	}

	return deactivatedAt
}

func (db *Database) RejoinedAt(ctx context.Context, user *tg.User) string {

	var (
		attempts   int
		retryDelay = constants.RetryDelay
		rejoinedAt string
	)

Retry:
	attempts++
	err := db.Postgres.QueryRowContext(ctx, `
		SELECT rejoined_at
		FROM user_configs
		WHERE id = $1;`,
		user.ID).Scan(&rejoinedAt)

	if err != nil {
		if strings.Contains(err.Error(), e.UnsupportedConversion) {
			return ""
		} else if strings.Contains(err.Error(), e.NotFound) {
			db.JoinUser(ctx, user, 0)
			goto Retry
		} else {
			log.Printf("can't get rejoined_at: %v", err)
			if attempts < constants.RetryAttempts {
				functions.Sleep(&retryDelay)
				goto Retry
			} else {
				log.Printf("getting rejoined_at failed after %d attempts", attempts)
			}
		}
	}

	return rejoinedAt
}
