package utils

import (
	"io"
	"log"
	"net/http"
	"os"
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

func Slice(s string, maxLen int) []string {
	var slices []string
	runes := []rune(s)
	for len(runes) > maxLen {
		cutIndex := maxLen
	Loop:
		for i := maxLen; i >= 0; i-- {
			switch runes[i] {
			case ' ', '\n', '\t', '\r':
				cutIndex = i
				break Loop
			}
		}
		slices = append(slices, string(runes[:cutIndex]))
		runes = runes[cutIndex:]
	}
	return append(slices, string(runes))
}

func MarkdownV2(s string) string {
	escapeChars := "\\_[]()~>#+-=|{}.!"
	for i := range escapeChars {
		char := string(escapeChars[i])
		escape := "\\" + char
		s = strings.ReplaceAll(s, char, escape)
	}

	if strings.Count(s, "```")%2 != 0 {
		s += "```"
	}

	var found1 bool
	var found2 bool
	var before1 string
	var before2 string
	var buffer strings.Builder

	for {
		before1, s, found1 = strings.Cut(s, "```")
		if strings.Count(before1, "**")%2 != 0 {
			before1 += "**"
		}

		for {
			before2, before1, found2 = strings.Cut(before1, "**")
			if strings.Count(before2, "`")%2 != 0 {
				before2 = strings.ReplaceAll(before2, "`", "\\`")
			}
			if strings.Count(before2, "*")%2 != 0 {
				before2 = strings.ReplaceAll(before2, "*", "\\*")
			}
			buffer.WriteString(before2)
			if !found2 {
				break
			}
			buffer.WriteString("*__")

			before2, before1, _ = strings.Cut(before1, "**")
			before2 = strings.ReplaceAll(before2, "`", "\\`")
			before2 = strings.ReplaceAll(before2, "*", "\\*")
			buffer.WriteString(before2)
			buffer.WriteString("__*")
		}
		if !found1 {
			break
		}
		buffer.WriteString("```")

		before1, s, _ = strings.Cut(s, "```")
		before1 = strings.ReplaceAll(before1, "`", "\\`")
		before1 = strings.ReplaceAll(before1, "*", "\\*")
		buffer.WriteString(before1)
		buffer.WriteString("```")
	}

	return buffer.String()
}

func DownloadFile(URL, path string) error {
	resp, err := http.Get(URL)
	if err != nil {
		log.Println("can't get file:", err)
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	out, err := os.Create(path)
	if err != nil {
		log.Println("can't create file:", err)
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println("can't write to file:", err)
		return err
	}
	return nil
}
