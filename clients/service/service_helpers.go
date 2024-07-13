package service

import (
	"context"
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
		s = item[0].(*regexp.Regexp).ReplaceAllString(s, item[1].(string))
	}
	r := []rune(s)
	var builder strings.Builder
	for i := 0; i < len(r); {
		if strings.HasPrefix(string(r[i:]), "\\frac{") && r[i+6] != '(' {
			r = append(r[:i+6], []rune(preProcess(string(r[i+6:])))...)
			for j := 0; j < 2; j++ {
				start, stack := -1, 0
				for ; i < len(r); i++ {
					if r[i] == '{' && r[i+1] != '(' {
						if stack == 0 {
							builder.WriteRune('{')
							start = i + 1
						}
						stack++
					} else if r[i] == '}' && stack != 0 {
						if stack == 1 {
							if strings.ContainsAny(string(r[start:i]), "+-*·×/÷^ ") {
								builder.WriteRune('(')
								builder.WriteString(string(r[start:i]))
								builder.WriteRune(')')
								builder.WriteRune('}')
							} else {
								builder.WriteString(string(r[start:i]))
								builder.WriteRune('}')
							}
							i++
							break
						}
						stack--
					} else if start == -1 {
						builder.WriteRune(r[i])
					}
				}
			}
		} else {
			builder.WriteRune(r[i])
			i++
		}
	}
	return builder.String()
}

func postProcess(s string) string {
	for _, item := range utils.PostProcessing {
		s = strings.ReplaceAll(s, item[0], item[1])
	}
	s = utils.SubSupRgx.ReplaceAllString(s, "($1; $2)")
	fractions := utils.FracRgx.FindAllString(s, -1)
	for _, original := range fractions {
		if len(original) > 20 {
			fraction := strings.ReplaceAll(original, "/", " / ")
			s = strings.Replace(s, original, fraction, 1)
		}
	}
	for strings.Contains(s, "  ") {
		s = strings.ReplaceAll(s, "  ", " ")
	}
	r := []rune(s)
	var builder, buffer strings.Builder
	for i := 0; i < len(r); i++ {
		if strings.HasPrefix(string(r[i:]), "_{") {
			start := i
			for i += 2; i < len(r); i++ {
				if r[i] == '}' {
					builder.WriteString(buffer.String())
					buffer.Reset()
					break
				}
				repl, ok := utils.Subscripts[r[i]]
				if !ok {
					builder.WriteString(string(r[start : i+1]))
					break
				}
				buffer.WriteString(repl)
			}
		} else if strings.HasPrefix(string(r[i:]), "^{") {
			start := i
			for i += 2; i < len(r); i++ {
				if r[i] == '}' {
					builder.WriteString(buffer.String())
					buffer.Reset()
					break
				}
				repl, ok := utils.Superscripts[r[i]]
				if !ok {
					builder.WriteString(string(r[start : i+1]))
					break
				}
				buffer.WriteString(repl)
			}
		} else {
			builder.WriteRune(r[i])
		}
	}
	return builder.String()
}
