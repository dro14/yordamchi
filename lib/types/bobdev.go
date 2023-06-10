package types

type TikToken struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Request struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	Price        int    `json:"price"`
	Requests     int    `json:"requests"`
}

type Response struct {
	Ok   bool   `json:"ok"`
	URL  string `json:"url"`
	Fake bool   `json:"fake"`
}
