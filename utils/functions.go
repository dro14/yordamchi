package utils

import (
	"log"
	"strings"
	"time"
)

func Sleep(retryDelay *time.Duration) {
	if *retryDelay > 0 {
		log.Printf("retrying request after %v", *retryDelay)
		time.Sleep(*retryDelay)
		*retryDelay *= 2
	}
}

func Slice(completion string, maxLen int) []string {
	var completions []string
	for len(completion) > maxLen {
		cutIndex := 0
	Loop:
		for i := maxLen; i >= 0; i-- {
			switch completion[i] {
			case ' ', '\n', '\t', '\r':
				cutIndex = i
				break Loop
			}
		}
		completions = append(completions, completion[:cutIndex])
		completion = completion[cutIndex:]
	}
	return append(completions, completion)
}

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
		if strings.Count(before, "`") > 0 {
			backticks := strings.Count(before, "`")
			shouldEscape := true
			if backticks%2 == 0 {
				shouldEscape = false
			}
			if shouldEscape {
				before = ReverseString(before)
				before = strings.Replace(before, "`", "\\`", 1)
				before = ReverseString(before)
			}
		}
		if strings.Count(before, "*") > 0 {
			doubleAsterisks := strings.Count(before, "**")
			shouldEscape := true
			if doubleAsterisks > 0 {
				before = strings.ReplaceAll(before, "**", Delim)
				if doubleAsterisks%2 != 0 {
					before = ReverseString(before)
					before = strings.Replace(before, Delim, "\\*\\_\\_", 1)
					before = ReverseString(before)
				}
				if strings.Count(before, "*")%2 != 0 {
					before = ReverseString(before)
					before = strings.Replace(before, "*", "\\*", 1)
					before = ReverseString(before)
				}
				isEnd := false
				for strings.Count(before, Delim) > 0 {
					if !isEnd {
						before = strings.Replace(before, Delim, "*__", 1)
					} else {
						before = strings.Replace(before, Delim, "__*", 1)
					}
					isEnd = !isEnd
				}
				shouldEscape = false
			}
			if shouldEscape {
				before = strings.ReplaceAll(before, "*", "\\*")
			}
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

func ReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
