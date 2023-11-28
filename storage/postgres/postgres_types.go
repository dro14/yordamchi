package postgres

type Message struct {
	ID               int
	UserID           int64
	Type             string
	CreatedOn        string
	PromptedAt       string
	CompletedAt      string
	FirstSend        int
	LastEdit         int
	PromptTokens     int
	PromptLength     int
	CompletionTokens int
	CompletionLength int
	Activity         int
	Requests         int
	Attempts         int
	FinishReason     string
	LanguageCode     string
}
