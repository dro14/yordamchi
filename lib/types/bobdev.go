package types

type TikToken struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}
