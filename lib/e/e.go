package e

import "errors"

// Telegram errors
var (
	UserDeletedMessage = errors.New("429 Bad Request: message to edit not found")
	UserBlockedBot     = errors.New("403 Forbidden: bot was blocked by the user")
)

// Postgres errors
const (
	UniqueConstraint      = "violates unique constraint"
	ForeignKeyConstraint  = "violates foreign key constraint"
	UnsupportedConversion = "converting NULL to string is unsupported"
	NotFound              = "no rows in result set"
)
