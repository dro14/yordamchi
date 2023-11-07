package bobdev

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dro14/yordamchi/lib/models"
	"github.com/dro14/yordamchi/lib/types"
)

type Request struct {
	Model    string          `json:"model"`
	Messages []types.Message `json:"messages"`
}

func Tokens(ctx context.Context, messages []types.Message) int {

	model := "gpt-3.5-turbo"
	if ctx.Value("model") == models.GPT4 {
		model = "gpt-4"
	}

	request := &Request{
		Model:    model,
		Messages: messages,
	}

	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(request)
	if err != nil {
		log.Printf("can't encode request: %v", err)
		return 0
	}

	resp, err := http.Post("https://yordamchi-d2hjy.ondigitalocean.app/fastapi-tiktoken2/tiktoken", "application/json", &buffer)
	if err != nil {
		log.Printf("can't send request: %v", err)
		return 0
	}

	response := make(map[string]int)
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0
	}
	_ = resp.Body.Close()

	return response["tokens"]
}
