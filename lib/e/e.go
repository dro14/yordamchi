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
	ResponseDecodeError = errors.New("can't decode response")
	UserDeletedMessage  = errors.New(MessageNotFound)
	UserBlockedError    = errors.New(UserBlocked)
)

// OpenAI errors
const (
	ContextLengthExceeded = "This model's maximum context length is 4097 tokens. However, you requested "
	InvalidRequest        = "invalid_request_error"
	StreamError           = "stream error"
	EmptyCompletion       = "empty completion"
	ServiceUnavailable    = "503 Service Unavailable"
	InternalServerError   = "500 Internal Server Error"
)

// Redis errors
const (
	KeyNotFound    = "redis: nil"
	UserNotDefined = "user is neither premium nor free"
)

var (
	UserNotDefinedError = errors.New(UserNotDefined)
)

// Postgres errors
const (
	UniqueConstraint      = "violates unique constraint"
	ForeignKeyConstraint  = "violates foreign key constraint"
	UnsupportedConversion = "converting NULL to string is unsupported"
	NotFound              = "no rows in result set"
)
