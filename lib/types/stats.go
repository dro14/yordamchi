package types

type Stats struct {
	IsPremium                          bool
	CompletedAt                        int64
	FirstSend, LastEdit                int64
	PromptTokens, PromptLength         int
	CompletionTokens, CompletionLength int
	Activity, Requests, Attempts       int
	FinishReason                       string
}
