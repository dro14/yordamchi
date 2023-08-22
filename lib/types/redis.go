package types

type UserStatus int

const (
	UnknownStatus UserStatus = iota
	BlockedStatus
	GPT4Status
	PremiumStatus
	FreeStatus
	ExhaustedStatus
)

type Activity struct {
	MessageID    int    `json:"message_id"`
	Message      string `json:"message"`
	Date         int    `json:"date"`
	UserID       int64  `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	TargetLang   string `json:"target_lang"`
	IsPremium    string `json:"is_premium"`
}
