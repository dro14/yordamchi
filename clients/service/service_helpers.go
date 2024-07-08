package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/dro14/yordamchi/utils"
)

func id(ctx context.Context) int64 {
	return ctx.Value("user_id").(int64)
}

func lang(ctx context.Context) string {
	return ctx.Value("language_code").(string)
}

func model(ctx context.Context) string {
	return ctx.Value("model").(string)
}

func preProcess(s string) string {
	for _, item := range utils.PreProcessing {
		s = regexp.MustCompile(item[0]).ReplaceAllString(s, item[1])
	}
	runes := []rune(s)
	i, builder := 0, strings.Builder{}
	for i < len(runes) {
		if strings.HasPrefix(string(runes[i:]), "\\frac{") && runes[i+6] != '(' {
			runes = append(runes[:i+6], []rune(preProcess(string(runes[i+6:])))...)
			for j := 0; j < 2; j++ {
				start, stack := -1, 0
				for ; i < len(runes); i++ {
					if runes[i] == '{' {
						if stack == 0 {
							builder.WriteRune('{')
							start = i + 1
						}
						stack++
					} else if runes[i] == '}' {
						if stack == 1 {
							if strings.ContainsAny(string(runes[start:i]), "+-*·×/÷^ ") {
								builder.WriteRune('(')
								builder.WriteString(string(runes[start:i]))
								builder.WriteRune(')')
								builder.WriteRune('}')
							} else {
								builder.WriteString(string(runes[start:i]))
								builder.WriteRune('}')
							}
							i++
							break
						}
						stack--
					} else if start == -1 {
						builder.WriteRune(runes[i])
					}
				}
			}
		} else {
			builder.WriteRune(runes[i])
			i++
		}
	}
	return builder.String()
}

func postProcess(s string) string {
	for _, item := range utils.PostProcessing {
		s = strings.ReplaceAll(s, item[0], item[1])
	}
	fractions := utils.FracRgx.FindAllString(s, -1)
	for _, original := range fractions {
		if len(original) > 20 {
			fraction := strings.ReplaceAll(original, "/", " / ")
			s = strings.Replace(s, original, fraction, 1)
		}
	}
	subSupers := utils.SubSuperRgx.FindAllStringSubmatch(s, -1)
	for _, match := range subSupers {
		subSuper := fmt.Sprintf("(%s, %s)", match[1], match[2])
		s = strings.Replace(s, match[0], subSuper, 1)
	}
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	return s
}
