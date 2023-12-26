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
	LaTeXRgx = regexp.MustCompile(`\\[(\[]\s?.+?\s?\\[)\]]`)
	TableRgx = regexp.MustCompile("(?m)(^```.*$\\s*)?(^\\|.*\\|$\\s*^\\|[-| :]*\\|$\\s*)(^\\|.*\\|$\\s*)*(^```$\\s*)?")
)

const (
	Fraction   = `\\\s?d?frac\s?{(.+?)}{(.+?)}`
	SquareRoot = `\\\s?sqrt\s?{(.+?)}`
)

var LaTeXReplacements = [][]string{
	// Greek letters
	{`\\alpha`, "α"},
	{`\\beta`, "β"},
	{`\\gamma`, "γ"},
	{`\\delta`, "δ"},
	{`\\(var)?epsilon`, "ε"},
	{`\\zeta`, "ζ"},
	{`\\eta`, "η"},
	{`\\(var)?(theta|teta)`, "θ"},
	{`\\iota`, "ι"},
	{`\\kappa`, "κ"},
	{`\\lambda`, "λ"},
	{`\\mu`, "μ"},
	{`\\nu`, "ν"},
	{`\\xi`, "ξ"},
	{`\\pi`, "π"},
	{`\\(var)?rho`, "ρ"},
	{`\\sigma`, "σ"},
	{`\\tau`, "τ"},
	{`\\upsilon`, "υ"},
	{`\\(var)?phi`, "φ"},
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
	{`\\rightarrow`, "→"},
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
	{`\\\s?(?:partial|qisman)\s?`, "∂"},
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

	{`\\\s?(?:sum|summa)\s?`, "Σ"},
	{`\\prod`, "Π"},
	{`\\binom`, "C"},
	{`\\int`, "∫"},
	{`\\iint`, "∬"},
	{`\\langle`, "⟨"},
	{`\\rangle`, "⟩"},
	{`\\pm`, "±"},
	{`\\mp`, "∓"},
	{`\\subset`, "⊂"},
	{`\\supset`, "⊃"},
	{`\\subseteq`, "⊆"},
	{`\\supseteq`, "⊇"},
	{`^\\circ\s?`, "°"},
	{`\\ldots`, "..."},
	{`\\\|`, "‖"},

	// Binary Operation/Relation Symbols
	{`\\\s?(?:times|marta|qat)`, "×"},
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
	{`\\\s?(?:approx|taxminan)`, "≈"},
	{`\\vee`, "∨"},
	{`\\otimes`, "⊗"},
	{`\\boxtimes`, "⊠"},
	{`\\cong`, "≅"},

	{`\\bar{r}`, "r̄"},
	{`\\bar{R}`, "R̄"},
	{`\\bar{x}`, "x̄"},
	{`\\bar{X}`, "X̄"},
	{`\\bar{y}`, "ȳ"},
	{`\\bar{Y}`, "Ȳ"},
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
	{`\\_`, "_"},
	{`\\ln`, "ln"},
	{`\\log`, "log"},
	{`\\sin`, "sin"},
	{`\\cos`, "cos"},
	{`\\tan`, "tan"},
	{`\\cot`, "cot"},
	{`\\sinh`, "sinh"},
	{`\\cosh`, "cosh"},
	{`\\tanh`, "tanh"},
	{`\\coth`, "coth"},
	{`\\arcsin`, "arcsin"},
	{`\\arccos`, "arccos"},
	{`\\arctan`, "arctan"},
	{`\\arccot`, "arccot"},
	{`\\arcsinh`, "arcsinh"},
	{`\\arccosh`, "arccosh"},
	{`\\arctanh`, "arctanh"},
	{`\\arccoth`, "arccoth"},

	{`\\\s?(?:text|matn|vec)\s?{(.+?)}`, "REPLACE"},
	{SquareRoot, "√(REPLACE)"},
	{Fraction, "(REPLACE)/(REPLACE)"},
	{`\\\s?(?:left|chap|right|o'ng|text|matn|limits|vec)\s?`, ""},
	{`\\(?: |,|;|:|quad)`, " "},
	{`\\[(\[]\s?(.+?)\s?\\[)\]]`, "`REPLACE`"},
}
