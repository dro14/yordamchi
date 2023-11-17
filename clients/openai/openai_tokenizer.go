package openai

import (
	"github.com/dro14/yordamchi/clients/openai/types"
)

func (o *OpenAI) countTokens(input any) int {
	switch input.(type) {
	case []types.Message:
		messages := input.([]types.Message)
		tokens := 0
		for i := range messages {
			content, ok := messages[i].Content.([]types.Content)
			if ok {
				if len(content) == 2 {
					messages[i].Content = content[1].Text
				} else {
					messages[i].Content = ""
				}
			}
			tokens += len(o.tkm.Encode(messages[i].Content.(string), nil, nil))
			tokens += len(o.tkm.Encode(messages[i].Role, nil, nil))
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
