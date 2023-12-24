package utils

import (
	"regexp"
	"time"
)

const (
	NumOfFreeReqs = 5
	RetryAttempts = 10
	RetryDelay    = 1000 * time.Millisecond
	ReqInterval   = 2000 * time.Millisecond
	Delim         = "\n-\n-\n-\n-\n"
)

var (
	LaTeXRgx = regexp.MustCompile(`\\[(|\[]\s?.+?\s?\\[)|\]]`)
	TableRgx = regexp.MustCompile("(?m)(^```.*$\\s*)?(^\\|.*\\|$\\s*^\\|[-| :]*\\|$\\s*)(^\\|.*\\|$\\s*)*(^```$\\s*)?")
)

var LaTeXReplacements = [][]string{
	{`\\(?:left|chap|right|o'ng|,|quad|text |limits)`, ""},

	{`\\(?:cdot|qat)`, "·"},
	{`\\(?:times|marta)`, "×"},
	{`\\(?:approx|taxminan)`, "≈"},
	{`\\pm`, "±"},
	{`\\mp`, "∓"},
	{`\\neq`, "≠"},
	{`\\leq`, "≤"},
	{`\\geq`, "≥"},
	{`\\cap`, "∩"},
	{`\\cup`, "∪"},
	{`\\subset`, "⊂"},
	{`\\supset`, "⊃"},
	{`\\subseteq`, "⊆"},
	{`\\supseteq`, "⊇"},
	{`\\(?:sum|summa)`, "Σ"},
	{`\\prod`, "Π"},
	{`\\infty`, "∞"},
	{`\\binom`, "C"},
	{`\\int`, "∫"},
	{`\\iint`, "∬"},
	{`\\(?:partial|qisman)`, "∂"},
	{`\\langle`, "⟨"},
	{`\\rangle`, "⟩"},
	{`^\\circ`, "°"},
	{`\\ldots`, "..."},
	{`\\limits`, ""},

	{`\\%`, "%"},
	{`\\ln`, "ln"},
	{`\\sin`, "sin"},
	{`\\cos`, "cos"},
	{`\\tan`, "tan"},
	{`\\cot`, "cot"},
	{`\\arcsin`, "arcsin"},
	{`\\arccos`, "arccos"},
	{`\\arctan`, "arctan"},
	{`\\arccot`, "arccot"},

	{`\\Delta`, "Δ"},
	{`\\Sigma`, "Σ"},
	{`\\Omega`, "Ω"},
	{`\\sigma`, "σ"},
	{`\\pi`, "π"},
	{`\\chi`, "χ"},
	{`\\omega`, "ω"},
	{`\\(?:theta|teta)`, "θ"},
	{`\\mu`, "μ"},
	{`\\phi`, "φ"},
	{`\\psi`, "ψ"},
	{`\\rho`, "ρ"},
	{`\\lambda`, "λ"},

	{`\\(?:text|matn|mathbb){(.+?)}`, "REPLACE"},
	{`\\sqrt{(.+?)}`, "√(REPLACE)"},
	{`\\frac{(.+?)}{(.+?)}`, "(REPLACE)/(REPLACE)"},
	{`\\[(|\[]\s?(.+?)\s?\\[)|\]]`, "`REPLACE`"},
}
