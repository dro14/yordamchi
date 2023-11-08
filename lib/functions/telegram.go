package functions

import (
	"strings"
)

func MarkdownV2(text string) string {

	escapeChars := "\\_[]()~>#+-=|{}.!"
	for i := range escapeChars {
		char := string(escapeChars[i])
		escape := "\\" + char
		text = strings.ReplaceAll(text, char, escape)
	}

	block := "```"
	blocks := strings.Count(text, block)
	if blocks%2 != 0 {
		text += block
	}

	var (
		found  bool
		before string
		buffer strings.Builder
	)

	for {
		before, text, found = strings.Cut(text, block)
		switch {
		case strings.Count(before, "`")%2 != 0:
			before = strings.ReplaceAll(before, "`", "\\`")
		case strings.Count(before, "*")%4 != 0 ||
			strings.Count(before, "**")%2 != 0:
			before = strings.ReplaceAll(before, "*", "\\*")
		}
		buffer.WriteString(before)
		if !found {
			break
		}
		buffer.WriteString(block)

		before, text, _ = strings.Cut(text, block)
		before = strings.ReplaceAll(before, "`", "\\`")
		before = strings.ReplaceAll(before, "*", "\\*")
		buffer.WriteString(before)
		buffer.WriteString(block)
	}

	return buffer.String()
}
