package service

import (
	"context"
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
		s = strings.ReplaceAll(s, item[0], item[1])
	}
	i, builder := 0, strings.Builder{}
	for i < len(s) {
		if strings.HasPrefix(s[i:], `\frac`) {
			for j := 0; j < 2; j++ {
				start, stack := -1, 0
				for ; i < len(s); i++ {
					if s[i] == '{' {
						if stack == 0 {
							builder.WriteString("{")
							start = i + 1
						}
						stack++
					} else if s[i] == '}' {
						if stack == 1 {
							if strings.ContainsAny(s[start:i], `+-*·×/÷^ `) {
								builder.WriteString("(" + s[start:i] + ")}")
							} else {
								builder.WriteString(s[start:i] + "}")
							}
							i++
							break
						}
						stack--
					} else if start == -1 {
						builder.WriteString(string(s[i]))
					}
				}
			}
		} else {
			builder.WriteString(string(s[i]))
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
			fraction := strings.Replace(original, "/", " / ", 1)
			s = strings.ReplaceAll(s, original, fraction)
		}
	}
	chars := []rune(s)
	i, builder, buffer := 0, strings.Builder{}, strings.Builder{}
	for i < len(chars) {
		if strings.HasPrefix(string(chars[i:]), `_{`) {
			start := i
			i += 2
			for ; i < len(chars); i++ {
				if chars[i] == '}' {
					builder.WriteString(buffer.String())
					buffer.Reset()
					i++
					break
				}
				char, ok := utils.Subscripts[chars[i]]
				if !ok {
					builder.WriteString(string(chars[start:i]))
					break
				}
				buffer.WriteString(char)
			}
		} else if strings.HasPrefix(string(chars[i:]), `^{`) {
			start := i
			i += 2
			for ; i < len(chars); i++ {
				if chars[i] == '}' {
					builder.WriteString(buffer.String())
					buffer.Reset()
					i++
					break
				}
				char, ok := utils.Superscripts[chars[i]]
				if !ok {
					builder.WriteString(string(chars[start:i]))
					break
				}
				buffer.WriteString(char)
			}
		} else {
			builder.WriteRune(chars[i])
			i++
		}
	}
	return "`" + builder.String() + "`"
}
