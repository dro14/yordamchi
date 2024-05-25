package openai

import (
	"github.com/dro14/yordamchi/clients/openai/types"
)

func (o *OpenAI) countTokens(input any) int {
	switch input.(type) {
	case []types.Message:
		tokens := 0
		for _, message := range input.([]types.Message) {
			content, ok := message.Content.([]types.Content)
			if ok {
				message.Content = content[0].Text
			}
			if message.Content != nil {
				tokens += len(o.tkm.Encode(message.Content.(string), nil, nil))
			}
			tokens += len(o.tkm.Encode(message.Role, nil, nil))
			tokens += 3
		}
		tokens += 3
		return tokens
	case string:
		return len(o.tkm.Encode(input.(string), nil, nil))
	default:
		return 0
	}
}
