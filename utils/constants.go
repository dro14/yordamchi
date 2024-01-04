package utils

import (
	"regexp"
	"time"
)

const (
	NumOfFreeReqs  = 5
	RetryAttempts  = 10
	RetryDelay     = 1000 * time.Millisecond
	ReqInterval    = 2000 * time.Millisecond
	NotifyInterval = 12 * time.Hour
	Delim          = "\n-\n-\n-\n-\n"
)

var (
	LaTeXRgx = regexp.MustCompile(`\\[(\[]\s?(.+?)\s?\\[)\]]`)
	TableRgx = regexp.MustCompile("(?m)(^```.*$\\s*)?(^\\|.*\\|$\\s*^\\|[-| :]*\\|$\\s*)(^\\|.*\\|$\\s*)*(^```$\\s*)?")
)

const (
	Fraction1 = `\\\s?d?frac\s?{(.+?)[})]\s?{(.+?)}`
	Fraction2 = `\\\s?d?frac\s?{(.+?){(.+?)}\s?}`
	Root2     = `\\\s?sqrt\s?{(.+?)}`
	Root3     = `\\\s?sqrt\s?[3]\s?{(.+?)}`
	Root4     = `\\\s?sqrt\s?[4]\s?{(.+?)}`
)

var LaTeXReplacements = [][]string{
	// Greek letters
	{`\\(?:alpha|alfa)`, "α"},
	{`\\beta`, "β"},
	{`\\gamma`, "γ"},
	{`\\delta`, "δ"},
	{`\\(?:var)?epsilon`, "ε"},
	{`\\zeta`, "ζ"},
	{`\\eta`, "η"},
	{`\\(?:var)?(theta|teta)`, "θ"},
	{`\\iota`, "ι"},
	{`\\kappa`, "κ"},
	{`\\lambda`, "λ"},
	{`\\mu`, "μ"},
	{`\\nu`, "ν"},
	{`\\xi`, "ξ"},
	{`\\pi`, "π"},
	{`\\(?:var)?rho`, "ρ"},
	{`\\sigma`, "σ"},
	{`\\tau`, "τ"},
	{`\\upsilon`, "υ"},
	{`\\(?:var)?phi`, "φ"},
	{`\\chi`, "χ"},
	{`\\psi`, "ψ"},
	{`\\omega`, "ω"},

	{`\\Gamma`, "Γ"},
	{`\\Delta`, "Δ"},
	{`\\Theta`, "Θ"},
	{`\\Lambda`, "Λ"},
	{`\\Xi`, "Ξ"},
	{`\\Pi`, "Π"},
	{`\\Sigma`, "Σ"},
	{`\\Upsilon`, "Υ"},
	{`\\Phi`, "Φ"},
	{`\\Psi`, "Ψ"},
	{`\\Omega`, "Ω"},

	// Arrows
	{`\\leftarrow`, "←"},
	{`\\(?:rightarrow|to)`, "→"},
	{`\\uparrow`, "↑"},
	{`\\downarrow`, "↓"},
	{`\\leftrightarrow`, "↔"},
	{`\\updownarrow`, "↕"},
	{`\\rightleftharpoons`, "⇌"},

	{`\\Leftarrow`, "⇐"},
	{`\\Rightarrow`, "⇒"},
	{`\\Uparrow`, "⇑"},
	{`\\Downarrow`, "⇓"},
	{`\\Leftrightarrow`, "⇔"},
	{`\\Updownarrow`, "⇕"},

	// Miscellaneous symbols
	{`\\infty`, "∞"},
	{`\\Re`, "ℜ"},
	{`\\nabla`, "∇"},
	{`\\(?:partial|qisman)`, "∂"},
	{`\\emptyset`, "∅"},
	{`\\wp`, "℘"},
	{`\\neg`, "¬"},
	{`\\square`, "□"},
	{`\\blacksquare`, "■"},
	{`\\forall`, "∀"},
	{`\\Im`, "ℑ"},
	{`\\exists`, "∃"},
	{`\\nexists`, "∄"},
	{`\\varnothing`, "∅"},
	{`\\complement`, "∁"},
	{`\\cdots`, "⋯"},
	{`\\surd`, "√"},
	{`\\triangle`, "△"},

	{`\\sum(?:ma)?`, "Σ"},
	{`\\prod`, "Π"},
	{`\\binom`, "C"},
	{`\\int`, "∫"},
	{`\\iint`, "∬"},
	{`\\iiint`, "∭"},
	{`\\langle`, "⟨"},
	{`\\rangle`, "⟩"},
	{`\\pm`, "±"},
	{`\\mp`, "∓"},
	{`\\subset`, "⊂"},
	{`\\supset`, "⊃"},
	{`\\subseteq`, "⊆"},
	{`\\supseteq`, "⊇"},
	{`^\\circ`, "°"},
	{`\\ldots`, "..."},
	{`\\\|`, "‖"},
	{`\\(?:Bigg?\||mid)`, "|"},

	// Binary Operation/Relation Symbols
	{`\\(?:times|kes)`, "×"},
	{`\\?marta`, "×"},
	{`\\div`, "÷"},
	{`\\cup`, "∪"},
	{`\\leq`, "≤"},
	{`\\in`, "∈"},
	{`\\notin`, "∉"},
	{`\\simeq`, "≃"},
	{`\\wedge`, "∧"},
	{`\\oplus`, "⊕"},
	{`\\Box`, "□"},
	{`\\equiv`, "≡"},
	{`\\cdot`, "·"},
	{`\\cap`, "∩"},
	{`\\neq`, "≠"},
	{`\\geq`, "≥"},
	{`\\perp`, "⊥"},
	{`\\(?:approx|taxminan)`, "≈"},
	{`\\vee`, "∨"},
	{`\\otimes`, "⊗"},
	{`\\boxtimes`, "⊠"},
	{`\\cong`, "≅"},
	{`\\land`, "∧"},
	{`\\lor`, "∨"},
	{`\\lnot?`, "¬"},
	{`\\lxor`, "⊻"},

	{`\\bar{x}`, "x̄"},
	{`\\bar{X}`, "X̄"},
	{`\\bar{y}`, "ȳ"},
	{`\\bar{Y}`, "Ȳ"},
	{`\\hat{p}`, "p̂"},
	{`\\hat{P}`, "P̂"},
	{`\\hat{r}`, "r̂"},
	{`\\hat{R}`, "R̂"},
	{`\\hat{x}`, "x̂"},
	{`\\hat{X}`, "X̂"},
	{`\\hat{y}`, "ŷ"},
	{`\\hat{Y}`, "Ŷ"},
	{`\\mathbb{N}`, "ℕ"},
	{`\\mathbb{Z}`, "ℤ"},
	{`\\mathbb{Q}`, "ℚ"},
	{`\\mathbb{R}`, "ℝ"},
	{`\\mathbb{C}`, "ℂ"},

	{`\\#`, "#"},
	{`\\$`, "$"},
	{`\\%`, "%"},
	{`\\&`, "&"},
	{`\\{`, "{"},
	{`\\}`, "}"},
	{`\\ln`, "ln"},
	{`\\log`, "log"},
	{`\\lim`, "lim"},
	{`\\sin`, "sin"},
	{`\\cos`, "cos"},
	{`\\tan`, "tan"},
	{`\\cot`, "cot"},
	{`\\sec`, "sec"},
	{`\\csc`, "csc"},
	{`\\arcsin`, "arcsin"},
	{`\\arccos`, "arccos"},
	{`\\arctan`, "arctan"},
	{`\\arccot`, "arccot"},
	{`\\arcsec`, "arcsec"},
	{`\\arccsc`, "arccsc"},

	{`\\\s?(?:text|matn)\s?{(.+?)}`, "REPLACE"},
	{Root2, "√(REPLACE)"},
	{Root3, "∛(REPLACE)"},
	{Root4, "∜(REPLACE)"},
	{Fraction1, "(REPLACE)/(REPLACE)"},
	{Fraction2, "(REPLACE)/(REPLACE)"},
	{`\\\s?(?:left|chap|right|o['ʻ]ng|text|matn)\s?`, ""},
	{`\\(?: |,|;|:|quad)`, " "},
}
