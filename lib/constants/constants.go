package constants

import "time"

const (
	ChatGPTURL  = "https://api.openai.com/v1/chat/completions"
	DonationURL = "https://payme.uz/60d6dbeb3632e1ceb8664de3"
	CheckoutURL = "https://chatgpt-payment.herokuapp.com/user_payme/"
	TokensURL   = "https://chatgpt-payment.herokuapp.com/tiktoken"
	AdminURL    = "https://t.me/yordamchiga_yordam"
)

const (
	RetryAttempts = 10
)

const (
	RetryDelay      = 10 * time.Second
	RequestInterval = 2 * time.Second
)
