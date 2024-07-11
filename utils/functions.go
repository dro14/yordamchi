package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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
	runes := []rune(Table(s))
	for len(runes) > maxLen {
		cutIndex := maxLen
	Loop:
		for i := maxLen; i >= 0; i-- {
			if runes[i] == '\n' {
				cutIndex = i
				break Loop
			}
		}
		slices = append(slices, string(runes[:cutIndex]))
		runes = runes[cutIndex:]
	}
	slices = append(slices, string(runes))
	for i := range slices {
		codeTags := PreRgx.FindAllString(slices[i], -1)
		if len(codeTags)%2 != 0 && i+1 < len(slices) {
			slices[i] = slices[i] + "```"
			slices[i+1] = codeTags[len(codeTags)-1] + slices[i+1]
		}
	}
	return slices
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

func MarkdownV2(s string) string {
	escapeChars := "\\_[]()~>#+-=|{}.!"
	for i := range escapeChars {
		s = strings.ReplaceAll(s, escapeChars[i:i+1], "\\"+escapeChars[i:i+1])
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
			matches := CodeRgx.FindAllString(before2, -1)
			for _, match := range matches {
				escaped := strings.ReplaceAll(match, "*", "\\*")
				before2 = strings.Replace(before2, match, escaped, 1)
			}
			matches = CodeRgx.Split(before2, -1)
			for _, match := range matches {
				if strings.Count(match, "*")%2 != 0 {
					escaped := strings.ReplaceAll(match, "*", "\\*")
					before2 = strings.Replace(before2, match, escaped, 1)
				}
			}
			matches = BoldRgx.FindAllString(before2, -1)
			for _, match := range matches {
				escaped := strings.ReplaceAll(match, "`", "\\`")
				before2 = strings.Replace(before2, match, escaped, 1)
			}
			matches = BoldRgx.Split(before2, -1)
			for _, match := range matches {
				if strings.Count(match, "`")%2 != 0 {
					escaped := strings.ReplaceAll(match, "`", "\\`")
					before2 = strings.Replace(before2, match, escaped, 1)
				}
			}
			buffer.WriteString(before2)
			if !found2 {
				break
			}
			buffer.WriteString("*_")

			before2, before1, _ = strings.Cut(before1, "**")
			before2 = strings.ReplaceAll(before2, "`", "\\`")
			before2 = strings.ReplaceAll(before2, "*", "\\*")
			buffer.WriteString(before2)
			buffer.WriteString("_*")
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
	s = buffer.String()
	parts := strings.Split(s, "```")
	for i := 0; i < len(parts); i += 2 {
		headers := HeaderRgx.FindAllString(parts[i], -1)
		for j, header := range headers {
			header = "__" + HeaderRgx.FindStringSubmatch(header)[1] + "__"
			parts[i] = strings.Replace(parts[i], headers[j], header, 1)
		}
		links := LinkRgx.FindAllString(parts[i], -1)
		for j, link := range links {
			name := LinkRgx.FindStringSubmatch(link)[1]
			link = LinkRgx.FindStringSubmatch(link)[2]
			link = strings.ReplaceAll(link, "\\", "")
			link = fmt.Sprintf("[%s](%s)", name, link)
			parts[i] = strings.Replace(parts[i], links[j], link, 1)
		}
	}
	return strings.Join(parts, "```")
}

func Table(input string) string {
	tablesIndexes := TableRgx.FindAllStringIndex(input, -1)
	start := 0
	var result strings.Builder
	for _, loc := range tablesIndexes {
		result.WriteString(input[start:loc[0]])

		lines := strings.Split(input[loc[0]:loc[1]], "\n")
		for i := 0; i < len(lines); i++ {
			lines[i] = strings.TrimSpace(lines[i])
			if !strings.HasPrefix(lines[i], "|") {
				lines = append(lines[:i], lines[i+1:]...)
				i--
			}
		}
		if len(lines) < 2 {
			log.Println("invalid input: the markdown table must have at least a header and a separator line")
			result.WriteString(input[loc[0]:loc[1]])
			continue
		}

		header := lines[0]
		headerCols := strings.Split(header, "|")
		columnWidths := make([]int, len(headerCols))
		for i, col := range headerCols {
			col = strings.TrimSpace(col)
			columnWidths[i] = len([]rune(col))
		}

		for _, line := range lines[2:] {
			rowCols := strings.Split(line, "|")
			for i, col := range rowCols {
				if i < len(columnWidths) {
					col = strings.TrimSpace(col)
					if columnWidths[i] < len([]rune(col)) {
						columnWidths[i] = len([]rune(col))
					}
				}
			}
		}

		result.WriteString("```\n")
		for lineIndex, line := range lines {
			if lineIndex != 1 {
				rowCols := strings.Split(line, "|")
				for i, col := range rowCols {
					col = strings.TrimSpace(col)
					if i > 0 && i < len(columnWidths)-1 {
						format := "%-" + strconv.Itoa(columnWidths[i]) + "s"
						if _, err := strconv.ParseFloat(col, 64); err == nil {
							format = "%" + strconv.Itoa(columnWidths[i]) + "s"
						}
						result.WriteString("| " + fmt.Sprintf(format, col) + " ")
					}
				}
			} else {
				for i := range columnWidths {
					if i > 0 && i < len(columnWidths)-1 {
						result.WriteString("|" + strings.Repeat("-", columnWidths[i]+2))
					}
				}
			}
			result.WriteString("|\n")
		}

		result.WriteString("```\n\n")
		start = loc[1]
	}

	result.WriteString(input[start:])
	return result.String()
}
