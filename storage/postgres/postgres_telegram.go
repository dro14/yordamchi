package postgres

import (
	"context"
	"log"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (p *Postgres) UserStarted(ctx context.Context, user *tgbotapi.User) {
	query := "INSERT INTO users (id, first_name, last_name, username, language_code, is_active, started_at) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING;"
	args := []any{user.ID, user.FirstName, user.LastName, user.UserName, lang(ctx), true, time.Now().Format(time.DateTime)}
	err := p.execTelegram(ctx, user, query, args)
	if err != nil {
		log.Printf("user %d: failed to join user: %s", user.ID, err)
		return
	}

	if !p.IsActive(ctx, user) {
		p.UserRestarted(ctx, user)
	}
}

func (p *Postgres) SaveMessage(ctx context.Context, user *tgbotapi.User, msg *Message) {
	query := "INSERT INTO messages (user_id, is_premium, created_on, prompted_at, completed_at, first_send, last_edit, prompt_tokens, prompt_length, completion_tokens, completion_length, activity, requests, attempts, finish_reason, language_code) VAlUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);"
	args := []any{user.ID, msg.IsPremium, msg.CreatedOn, msg.PromptedAt, msg.CompletedAt, msg.FirstSend, msg.LastEdit, msg.PromptTokens, msg.PromptLength, msg.CompletionTokens, msg.CompletionLength, msg.Activity, msg.Requests, msg.Attempts, msg.FinishReason, lang(ctx)}
	err := p.execTelegram(ctx, user, query, args)
	if err != nil {
		log.Printf("user %d: failed to save message: %s", user.ID, err)
	}

	if !p.IsActive(ctx, user) {
		p.UserRestarted(ctx, user)
	}
}

func (p *Postgres) UserBlocked(ctx context.Context, user *tgbotapi.User) {
	query := "UPDATE users SET is_active = $1, blocked_at = $2 WHERE id = $3;"
	args := []any{false, time.Now().Format(time.DateTime), user.ID}
	err := p.execTelegram(ctx, user, query, args)
	if err != nil {
		log.Printf("user %d: failed to deactivate user: %s", user.ID, err)
	}
}

func (p *Postgres) UserRestarted(ctx context.Context, user *tgbotapi.User) {
	query := "UPDATE users SET is_active = $1, restarted_at = $2 WHERE id = $3;"
	args := []any{true, time.Now().Format(time.DateTime), user.ID}
	err := p.execTelegram(ctx, user, query, args)
	if err != nil {
		log.Printf("user %d: failed to rejoin user: %s", user.ID, err)
	}
}

func (p *Postgres) IsActive(ctx context.Context, user *tgbotapi.User) bool {
	query := "SELECT is_active FROM users WHERE id = $1;"
	args := []any{user.ID}
	var isActive bool
	err := p.queryTelegram(ctx, user, query, args, &isActive)
	if err != nil {
		log.Printf("user %d: failed to get is_active: %s", user.ID, err)
		return false
	}
	return isActive
}

func (p *Postgres) SetLang(ctx context.Context, user *tgbotapi.User) {
	query := "UPDATE users SET language_code = $1 WHERE id = $2;"
	args := []any{lang(ctx), user.ID}
	err := p.execTelegram(ctx, user, query, args)
	if err != nil {
		log.Printf("user %d: failed to set language_code: %s", user.ID, err)
	}
}
