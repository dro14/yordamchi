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
	{`\\(?:left|chap|right|o'ng|,|quad|text\s?|limits)`, ""},

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

	// Binary Operation/Relation Symbols
	{`\\(?:times|\s?marta)`, "×"},
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
	{`\\(?:cdot|qat)`, "·"},
	{`\\cap`, "∩"},
	{`\\neq`, "≠"},
	{`\\geq`, "≥"},
	{`\\perp`, "⊥"},
	{`\\(?:approx|taxminan)`, "≈"},
	{`\\vee`, "∨"},
	{`\\otimes`, "⊗"},
	{`\\boxtimes`, "⊠"},
	{`\\cong`, "≅"},

	{`\\pm`, "±"},
	{`\\mp`, "∓"},
	{`\\subset`, "⊂"},
	{`\\supset`, "⊃"},
	{`\\subseteq`, "⊆"},
	{`\\supseteq`, "⊇"},
	{`\\(?:sum|summa)\s?`, "Σ"},
	{`\\prod\s?`, "Π"},
	{`\\binom`, "C"},
	{`\\int`, "∫"},
	{`\\iint`, "∬"},
	{`\\langle`, "⟨"},
	{`\\rangle`, "⟩"},
	{`^\\circ\s?`, "°"},
	{`\\ldots`, "..."},
	{`\\bar{x}`, "x̄"},
	{`\\bar{X}`, "X̄"},
	{`\\bar{y}`, "ȳ"},
	{`\\bar{Y}`, "Ȳ"},
	{`\\hat{x}`, "x̂"},
	{`\\hat{X}`, "X̂"},
	{`\\hat{y}`, "ŷ"},
	{`\\hat{Y}`, "Ŷ"},
	{`\\mathbb{N}`, "ℕ"},
	{`\\mathbb{Z}`, "ℤ"},
	{`\\mathbb{Q}`, "ℚ"},
	{`\\mathbb{R}`, "ℝ"},
	{`\\mathbb{C}`, "ℂ"},

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

	{`\\(?:text|matn){(.+?)}`, "REPLACE"},
	{`\\sqrt{(.+?)}`, "√(REPLACE)"},
	{`\\frac{(.+?)}{(.+?)}`, "(REPLACE)/(REPLACE)"},
	{`\\[(|\[]\s?(.+?)\s?\\[)|\]]`, "`REPLACE`"},
}
