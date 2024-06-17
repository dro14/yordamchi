package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dro14/yordamchi/storage/postgres/types"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Postgres) UserStarted(ctx context.Context, user *tgbotapi.User) {
	query := "INSERT INTO users (id, first_name, last_name, username, language_code, is_active, started_at) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING;"
	args := []any{user.ID, user.FirstName, user.LastName, user.UserName, lang(ctx), true, time.Now().Format(time.DateTime)}
	err := p.execTelegram(ctx, user, query, args)
	if err != nil {
		log.Printf("user %d: failed to join user: %s", user.ID, err)
	}
}

func (p *Postgres) SaveMessage(ctx context.Context, user *tgbotapi.User, msg *types.Message) {
	query := "INSERT INTO messages (user_id, type, created_on, prompted_at, completed_at, first_send, last_edit, prompt_tokens, prompt_length, completion_tokens, completion_length, activity, requests, attempts, finish_reason, language_code, input, output) VAlUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18);"
	args := []any{user.ID, msg.Type, msg.CreatedOn, msg.PromptedAt, msg.CompletedAt, msg.FirstSend, msg.LastEdit, msg.PromptTokens, msg.PromptLength, msg.CompletionTokens, msg.CompletionLength, msg.Activity, msg.Requests, msg.Attempts, msg.FinishReason, lang(ctx), msg.Input, msg.Output}
	err := p.execTelegram(ctx, user, query, args)
	if err != nil {
		log.Printf("user %d: failed to save message: %s", user.ID, err)
	}
}

func (p *Postgres) UserBlocked(ctx context.Context, user *tgbotapi.User) {
	query := "UPDATE users SET first_name = $1, last_name = $2, username = $3, language_code = $4, is_active = $5, blocked_at = $6 WHERE id = $7;"
	args := []any{user.FirstName, user.LastName, user.UserName, lang(ctx), false, time.Now().Format(time.DateTime), user.ID}
	err := p.execTelegram(ctx, user, query, args)
	if err != nil {
		log.Printf("user %d: failed to deactivate user: %s", user.ID, err)
	}
}

func (p *Postgres) UserRestarted(ctx context.Context, user *tgbotapi.User) {
	query := "UPDATE users SET first_name = $1, last_name = $2, username = $3, language_code = $4, is_active = $5, restarted_at = $6 WHERE id = $7;"
	args := []any{user.FirstName, user.LastName, user.UserName, lang(ctx), true, time.Now().Format(time.DateTime), user.ID}
	err := p.execTelegram(ctx, user, query, args)
	if err != nil {
		log.Printf("user %d: failed to rejoin user: %s", user.ID, err)
	}
}

func (p *Postgres) SetLang(ctx context.Context, user *tgbotapi.User) {
	query := "UPDATE users SET language_code = $1 WHERE id = $2;"
	args := []any{lang(ctx), user.ID}
	err := p.execTelegram(ctx, user, query, args)
	if err != nil {
		log.Printf("user %d: failed to set language_code: %s", user.ID, err)
	}
}

func (p *Postgres) UpdateUser(ctx context.Context, user *tgbotapi.User) {
	query := "UPDATE users SET first_name = $1, last_name = $2, username = $3, language_code = $4, is_active = true WHERE id = $5;"
	args := []any{user.FirstName, user.LastName, user.UserName, lang(ctx), user.ID}
	err := p.execTelegram(ctx, user, query, args)
	if err != nil {
		log.Printf("user %d: failed to update user: %s", user.ID, err)
	}
}

func (p *Postgres) PollAnswer(ctx context.Context, pollAnswer *tgbotapi.PollAnswer) {
	if len(pollAnswer.OptionIDs) != 1 {
		log.Printf("user %d: poll_answer.option_ids length is not 1", pollAnswer.User.ID)
		return
	}
	query := "INSERT INTO poll_answers (poll_id, question, user_id, option_id) VALUES ($1, $2, $3, $4);"
	args := []any{pollAnswer.PollID, p.redis.PollQuestion(ctx), pollAnswer.User.ID, pollAnswer.OptionIDs[0]}
	err := p.execTelegram(ctx, &pollAnswer.User, query, args)
	if err != nil {
		log.Printf("user %d: failed to save poll answer: %s", pollAnswer.User.ID, err)
	}
}

func (p *Postgres) User(ctx context.Context, user *tgbotapi.User) string {
	query := "SELECT first_name, last_name, username FROM users WHERE id = $1;"
	var firstName, lastName, username string
	args := []any{user.ID}
	err := p.queryTelegram(ctx, user, query, args, &firstName, &lastName, &username)
	if err != nil {
		log.Printf("user %d: failed to get user: %s", user.ID, err)
		return ""
	}

	fullName := strings.TrimSpace(fmt.Sprintf("%s %s", firstName, lastName))
	if fullName != "" {
		return fullName
	}
	if username != "" {
		return "@" + username
	}
	return ""
}
