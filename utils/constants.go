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
	Fraction1 = `\\d?frac{(.+?)}{(.+?)}`
	Fraction2 = `\\d?frac{(.+?){(.+?)}}`
	Root2     = `\\sqrt{(.+?)}`
	Root3     = `\\sqrt[3]{(.+?)}`
	Root4     = `\\sqrt[4]{(.+?)}`
)

var LaTeXReplacements = [][]string{
	// Greek letters
	{`\\alpha`, "α"},
	{`\\beta`, "β"},
	{`\\gamma`, "γ"},
	{`\\delta`, "δ"},
	{`\\(?:var)?epsilon`, "ε"},
	{`\\zeta`, "ζ"},
	{`\\eta`, "η"},
	{`\\(?:var)?theta`, "θ"},
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
	{`\\partial`, "∂"},
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

	{`\\sum`, "Σ"},
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
	{`\\land`, "∧"},
	{`\\lor`, "∨"},
	{`\\lnot?`, "¬"},
	{`\\lxor`, "⊻"},
	{`\\lfloor`, "⌊"},
	{`\\rfloor`, "⌋"},
	{`\\lceil`, "⌈"},
	{`\\rceil`, "⌉"},

	// Binary Operation/Relation Symbols
	{`\\times`, "×"},
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
	{`\\approx`, "≈"},
	{`\\vee`, "∨"},
	{`\\otimes`, "⊗"},
	{`\\boxtimes`, "⊠"},
	{`\\cong`, "≅"},

	{`\\(?:bar|vec){x}`, "x̄"},
	{`\\(?:bar|vec){X}`, "X̄"},
	{`\\(?:bar|vec){y}`, "ȳ"},
	{`\\(?:bar|vec){Y}`, "Ȳ"},
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
	{`\^0`, "⁰"},
	{`\^1`, "¹"},
	{`\^2`, "²"},
	{`\^3`, "³"},
	{`\^4`, "⁴"},
	{`\^5`, "⁵"},
	{`\^6`, "⁶"},
	{`\^7`, "⁷"},
	{`\^8`, "⁸"},
	{`\^9`, "⁹"},
	{`\^n`, "ⁿ"},
	{`\^x`, "ˣ"},

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

	{`\\(?:text|mathbf|mathcal){(.+?)}`, "REPLACE"},
	{Root2, "√(REPLACE)"},
	{Root3, "∛(REPLACE)"},
	{Root4, "∜(REPLACE)"},
	{Fraction1, "(REPLACE)/(REPLACE)"},
	{Fraction2, "(REPLACE)/(REPLACE)"},
	{`\\(?:left|right)`, ""},
	{`\\(?: |,|;|:|quad)`, " "},
}
