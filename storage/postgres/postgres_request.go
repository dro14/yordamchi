package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dro14/yordamchi/utils"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Postgres) execTelegram(ctx context.Context, user *tgbotapi.User, query string, args []any) error {
	retryDelay := utils.RetryDelay
	attempts := 0
Retry:
	attempts++
	_, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("%s\nuser %d: %s\n", query, user.ID, err)
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "violates foreign key constraint") {
			p.UserStarted(ctx, user)
			goto Retry
		} else if attempts < utils.RetryAttempts {
			utils.Sleep(&retryDelay)
			goto Retry
		} else {
			return err
		}
	}
	return nil
}

func (p *Postgres) queryTelegram(ctx context.Context, user *tgbotapi.User, query string, args []any, results ...any) error {
	retryDelay := utils.RetryDelay
	attempts := 0
Retry:
	attempts++
	err := p.db.QueryRowContext(ctx, query, args...).Scan(results...)
	if err != nil {
		log.Printf("%s\nuser %d: %s\n", query, user.ID, err)
		if errors.Is(err, sql.ErrNoRows) {
			p.UserStarted(ctx, user)
			goto Retry
		} else if attempts < utils.RetryAttempts {
			utils.Sleep(&retryDelay)
			goto Retry
		} else {
			return err
		}
	}
	return nil
}

func (p *Postgres) execPayme(query string, args []any) error {
	retryDelay := utils.RetryDelay
	attempts := 0
Retry:
	attempts++
	_, err := p.db.Exec(query, args...)
	if err != nil {
		log.Printf("%s\n%s\n", query, err)
		if errors.Is(err, sql.ErrNoRows) {
			return err
		} else if attempts < utils.RetryAttempts {
			utils.Sleep(&retryDelay)
			goto Retry
		} else {
			return err
		}
	}
	return nil
}

func (p *Postgres) queryPayme(query string, args []any, results ...any) error {
	retryDelay := utils.RetryDelay
	attempts := 0
Retry:
	attempts++
	err := p.db.QueryRow(query, args...).Scan(results...)
	if err != nil {
		log.Printf("%s\n%s\n", query, err)
		if errors.Is(err, sql.ErrNoRows) {
			return err
		} else if attempts < utils.RetryAttempts {
			utils.Sleep(&retryDelay)
			goto Retry
		} else {
			return err
		}
	}
	return nil
}
