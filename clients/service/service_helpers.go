package service

import (
	"context"
	"strings"
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

var preProcessing = [][]string{
	{`\cot`, `cot`},
	{`\cross`, `×`},
	{`\implies`, `⇒`},
	{`\times`, `×`},
	{`\cdot`, `·`},
	{`\div`, `÷`},
}

func preProcess(s string) string {
	for _, item := range preProcessing {
		s = strings.ReplaceAll(s, item[0], item[1])
	}
	i, temp := 0, ""
	for i < len(s) {
		if strings.HasPrefix(s, `\frac`) {
			for j := 0; j < 2; j++ {
				start, stack := -1, 0
				for ; i < len(s); i++ {
					if s[i] == '{' {
						if stack == 0 {
							temp += "{"
							start = i + 1
						}
						stack++
					} else if s[i] == '}' {
						if stack == 1 {
							if strings.ContainsAny(s[start:i], `+-*·×/÷^ `) {
								temp += "(" + s[start:i] + ")}"
							} else {
								temp += s[start:i] + "}"
							}
							i++
							break
						}
						stack--
					} else if start == -1 {
						temp += string(s[i])
					}
				}
			}
		} else {
			temp += string(s[i])
			i++
		}
	}
	return s
}

var postProcessing = [][]string{

	// Subscripts
	{`_0`, `₀`},
	{`_1`, `₁`},
	{`_2`, `₂`},
	{`_3`, `₃`},
	{`_4`, `₄`},
	{`_5`, `₅`},
	{`_6`, `₆`},
	{`_7`, `₇`},
	{`_8`, `₈`},
	{`_9`, `₉`},
	{`_a`, `ₐ`},
	{`_e`, `ₑ`},
	{`_h`, `ₕ`},
	{`_i`, `ᵢ`},
	{`_j`, `ⱼ`},
	{`_k`, `ₖ`},
	{`_l`, `ₗ`},
	{`_m`, `ₘ`},
	{`_n`, `ₙ`},
	{`_o`, `ₒ`},
	{`_p`, `ₚ`},
	{`_r`, `ᵣ`},
	{`_s`, `ₛ`},
	{`_t`, `ₜ`},
	{`_u`, `ᵤ`},
	{`_v`, `ᵥ`},
	{`_x`, `ₓ`},
	{`_y`, `ᵧ`},
	{`_-`, `₋`},
	{`_+`, `₊`},
	{`_=`, `₌`},
	{`_*`, `⁎`},
	{`_π`, `ₚᵢ`},

	// Superscripts
	{`^0`, `⁰`},
	{`^1`, `¹`},
	{`^2`, `²`},
	{`^3`, `³`},
	{`^4`, `⁴`},
	{`^5`, `⁵`},
	{`^6`, `⁶`},
	{`^7`, `⁷`},
	{`^8`, `⁸`},
	{`^9`, `⁹`},
	{`^a`, `ᵃ`},
	{`^b`, `ᵇ`},
	{`^c`, `ᶜ`},
	{`^d`, `ᵈ`},
	{`^e`, `ᵉ`},
	{`^f`, `ᶠ`},
	{`^g`, `ᵍ`},
	{`^h`, `ʰ`},
	{`^i`, `ⁱ`},
	{`^j`, `ʲ`},
	{`^k`, `ᵏ`},
	{`^l`, `ˡ`},
	{`^m`, `ᵐ`},
	{`^n`, `ⁿ`},
	{`^o`, `ᵒ`},
	{`^p`, `ᵖ`},
	{`^r`, `ʳ`},
	{`^s`, `ˢ`},
	{`^t`, `ᵗ`},
	{`^u`, `ᵘ`},
	{`^v`, `ᵛ`},
	{`^w`, `ʷ`},
	{`^x`, `ˣ`},
	{`^y`, `ʸ`},
	{`^z`, `ᶻ`},
	{`^A`, `ᴬ`},
	{`^B`, `ᴮ`},
	{`^D`, `ᴰ`},
	{`^E`, `ᴱ`},
	{`^G`, `ᴳ`},
	{`^H`, `ᴴ`},
	{`^I`, `ᴵ`},
	{`^J`, `ᴶ`},
	{`^K`, `ᴷ`},
	{`^L`, `ᴸ`},
	{`^M`, `ᴹ`},
	{`^N`, `ᴺ`},
	{`^O`, `ᴼ`},
	{`^P`, `ᴾ`},
	{`^R`, `ᴿ`},
	{`^T`, `ᵀ`},
	{`^U`, `ᵁ`},
	{`^V`, `ⱽ`},
	{`^W`, `ᵂ`},
	{`^-`, `⁻`},
	{`^+`, `⁺`},
	{`^=`, `⁼`},
	{`^*`, `ˣ`},
	{`^π`, `ᵖⁱ`},
	{`^∘`, `°`},
}

func postProcess(s string) string {
	s = "`" + s + "`"
	for _, item := range postProcessing {
		s = strings.ReplaceAll(s, item[0], item[1])
	}
	return s
}
