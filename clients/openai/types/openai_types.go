package types

// Completions is a struct for OpenAI Completions API
type Completions struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Tools       []Tool    `json:"tools,omitempty"`
	User        string    `json:"user,omitempty"`
}

type Message struct {
	Role       string     `json:"role"`
	Content    any        `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
	Arguments   string         `json:"arguments,omitempty"`
}

type Content struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL           string `json:"url"`
	Detail        string `json:"detail,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

// Generations is a struct for OpenAI Generations API
type Generations struct {
	Prompt  string `json:"prompt"`
	Model   string `json:"model"`
	Quality string `json:"quality"`
	Size    string `json:"size"`
	Style   string `json:"style"`
	User    string `json:"user"`
}

// Response is a struct for OpenAI API response
type Response struct {
	Choices []Choice   `json:"choices"`
	Error   Error      `json:"error"`
	Data    []ImageURL `json:"data"`
}

type Choice struct {
	Message       Message       `json:"message"`
	Delta         Delta         `json:"delta"`
	FinishReason  string        `json:"finish_reason"`
	FinishDetails FinishDetails `json:"finish_details"`
}

type Delta struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls"`
}

type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type FinishDetails struct {
	Type string `json:"type"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}
