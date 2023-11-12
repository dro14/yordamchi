package e

import "errors"

// Telegram errors
var (
	UserDeletedMessage = errors.New("user deleted message")
	Forbidden          = errors.New("forbidden")
)

// Postgres errors
const (
	UniqueConstraint      = "violates unique constraint"
	ForeignKeyConstraint  = "violates foreign key constraint"
	UnsupportedConversion = "converting NULL to string is unsupported"
	NotFound              = "no rows in result set"
)
