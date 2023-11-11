package types

type Completions struct {
	Model     string    `json:"model,omitempty"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
	Stream    bool      `json:"stream,omitempty"`
	User      string    `json:"user,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type Content struct {
	Type     string   `json:"type"`
	Text     string   `json:"text,omitempty"`
	ImageURL ImageURL `json:"image_url,omitempty"`
}

type Response struct {
	Choices []Choice   `json:"choices"`
	Usage   Usage      `json:"usage"`
	Error   Error      `json:"error"`
	Data    []ImageURL `json:"data"`
}

type Choice struct {
	Message      Message `json:"message"`
	Delta        Delta   `json:"delta"`
	FinishReason string  `json:"finish_reason"`
}

type Delta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type Generations struct {
	Prompt  string `json:"prompt"`
	Model   string `json:"model"`
	Quality string `json:"quality"`
	Size    string `json:"size"`
	Style   string `json:"style"`
	User    string `json:"user"`
}
