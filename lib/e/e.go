package e

import "errors"

// Telegram errors
const (
	BrokenPipe         = "broken pipe"
	TooManyRequests    = "420: FLOOD_WAIT "
	MessageNotFound    = "400: MESSAGE_ID_INVALID"
	MessageEmpty       = "400: MESSAGE_EMPTY"
	MessageNotModified = "400: MESSAGE_NOT_MODIFIED"
	MessageTooLong     = "400: MESSAGE_TOO_LONG"
	UserBlocked        = "400: USER_IS_BLOCKED"
)

var (
	UserDeletedMessage = errors.New(MessageNotFound)
	UserBlockedError   = errors.New(UserBlocked)
)

const StreamError = "stream error" // OpenAI stream error

// Postgres errors
const (
	UniqueConstraint      = "violates unique constraint"
	ForeignKeyConstraint  = "violates foreign key constraint"
	UnsupportedConversion = "converting NULL to string is unsupported"
	NotFound              = "no rows in result set"
)
