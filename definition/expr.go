package definition

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type (
	Expr interface {
		fmt.Stringer
		exprPrecedence() int
		Children() []Expr
	}
	ExprCore interface {
		Expr
		exprCore()
	}
	Terminals interface {
		ExprCore
		isTerminal()
	}
	TextTerminals interface {
		Terminals
		Match([]byte) (int, bool)
	}

	Atom struct {
		Symbol       string
		Attributes   map[string][]byte
		Text         []byte
		TextSelector Segments
	}

	Choice     struct{ Exprs []Expr }
	Junction   struct{ Exprs []Expr }
	Negation   struct{ Expr Expr }
	Ensure     struct{ Expr Expr }
	Optional   struct{ Expr Expr }
	Kleene     struct{ Expr Expr }
	Repetition struct {
		Expr Expr
		Min  uint
	}
	Symbol struct {
		Name       string
		Attributes map[string][]byte
	}
	Empty       struct{}
	Dot         struct{}
	TextToken   struct{ Text []byte }
	AtomPattern struct{ Matcher map[string]TextTerminals }
	TextPattern struct {
		Expr  string
		Regex *regexp.Regexp
	}
	StartOfFile struct{}
	EndOfFile   struct{}
)

func (a Atom) SelectString() string {
	return string(a.SelectText())
}

func (a Atom) SelectText() []byte {
	if len(a.TextSelector.segments) == 0 {
		return nil
	}
	if len(a.TextSelector.segments) == 1 {
		return a.Text[a.TextSelector.segments[0].Start:a.TextSelector.segments[0].End]
	}
	totalLength := 0
	for _, selector := range a.TextSelector.segments {
		totalLength += selector.End - selector.Start
	}
	text := make([]byte, totalLength)
	offset := 0
	for _, selector := range a.TextSelector.segments {
		length := selector.End - selector.Start
		copy(text[offset:offset+length], a.Text[selector.Start:selector.End])
		offset += length
	}
	return text
}

func joinExprs(rootPrecedence int, items []Expr, separator string) string {
	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, wrapExpr(rootPrecedence, item))
	}
	return strings.Join(values, separator)
}

func wrapExpr(rootPrecedence int, item Expr) string {
	if item.exprPrecedence() < rootPrecedence {
		return "(" + item.String() + ")"
	}
	return item.String()
}

const (
	StartOfFileBuiltinSymbol = "@sof"
	EndOfFileBuiltinSymbol   = "@eof"
)

func (e Choice) String() string   { return joinExprs(e.exprPrecedence(), e.Exprs, " / ") }
func (e Junction) String() string { return joinExprs(e.exprPrecedence(), e.Exprs, " ") }
func (e Negation) String() string { return "!" + wrapExpr(e.exprPrecedence(), e.Expr) }
func (e Ensure) String() string   { return "&" + wrapExpr(e.exprPrecedence(), e.Expr) }
func (e Optional) String() string { return wrapExpr(e.exprPrecedence(), e.Expr) + "?" }
func (e Kleene) String() string   { return wrapExpr(e.exprPrecedence(), e.Expr) + "*" }
func (e Repetition) String() string {
	var suffix string
	if e.Min == 0 {
		suffix = "*"
	} else if e.Min == 1 {
		suffix = "+"
	} else {
		suffix = fmt.Sprintf("{%v,}", e.Min)
	}
	return wrapExpr(e.exprPrecedence(), e.Expr) + suffix
}
func (e Symbol) String() string      { return e.Name }
func (e Empty) String() string       { return "@empty" }
func (e Dot) String() string         { return "." }
func (e StartOfFile) String() string { return StartOfFileBuiltinSymbol }
func (e EndOfFile) String() string   { return EndOfFileBuiltinSymbol }
func (e TextPattern) String() string { return "=~" + strconv.Quote(e.Expr) }
func (e TextToken) String() string   { return strconv.Quote(string(e.Text)) }
func (e AtomPattern) String() string {
	attributes := make([]string, 0)
	for attributeKey, attributeMatcher := range e.Matcher {
		if attributeMatcher == nil {
			attributes = append(attributes, attributeKey)
		} else {
			attributes = append(attributes, fmt.Sprintf("%v:%v", attributeKey, attributeMatcher))
		}
	}
	return fmt.Sprintf("{%v}", strings.Join(attributes, ", "))
}

func (e Choice) exprPrecedence() int      { return 1 }
func (e Junction) exprPrecedence() int    { return 2 }
func (e Negation) exprPrecedence() int    { return 3 }
func (e Ensure) exprPrecedence() int      { return 3 }
func (e Optional) exprPrecedence() int    { return 4 }
func (e Kleene) exprPrecedence() int      { return 4 }
func (e Repetition) exprPrecedence() int  { return 4 }
func (e Symbol) exprPrecedence() int      { return 5 }
func (e Empty) exprPrecedence() int       { return 5 }
func (e Dot) exprPrecedence() int         { return 5 }
func (e TextToken) exprPrecedence() int   { return 5 }
func (e TextPattern) exprPrecedence() int { return 5 }
func (e AtomPattern) exprPrecedence() int { return 5 }
func (e StartOfFile) exprPrecedence() int { return 5 }
func (e EndOfFile) exprPrecedence() int   { return 5 }

func (e Choice) Children() []Expr      { return e.Exprs }
func (e Junction) Children() []Expr    { return e.Exprs }
func (e Negation) Children() []Expr    { return []Expr{e.Expr} }
func (e Ensure) Children() []Expr      { return []Expr{e.Expr} }
func (e Optional) Children() []Expr    { return []Expr{e.Expr} }
func (e Kleene) Children() []Expr      { return []Expr{e.Expr} }
func (e Repetition) Children() []Expr  { return []Expr{e.Expr} }
func (e Symbol) Children() []Expr      { return nil }
func (e Empty) Children() []Expr       { return nil }
func (e Dot) Children() []Expr         { return nil }
func (e TextToken) Children() []Expr   { return nil }
func (e TextPattern) Children() []Expr { return nil }
func (e AtomPattern) Children() []Expr { return nil }
func (e StartOfFile) Children() []Expr { return nil }
func (e EndOfFile) Children() []Expr   { return nil }

func (e Choice) exprCore()      {}
func (e Junction) exprCore()    {}
func (e Negation) exprCore()    {}
func (e Kleene) exprCore()      {}
func (e Symbol) exprCore()      {}
func (e Empty) exprCore()       {}
func (e Dot) exprCore()         {}
func (e TextToken) exprCore()   {}
func (e TextPattern) exprCore() {}
func (e AtomPattern) exprCore() {}
func (e StartOfFile) exprCore() {}
func (e EndOfFile) exprCore()   {}

func (e Empty) isTerminal()       {}
func (e Dot) isTerminal()         {}
func (e TextToken) isTerminal()   {}
func (e TextPattern) isTerminal() {}
func (e AtomPattern) isTerminal() {}
func (e StartOfFile) isTerminal() {}
func (e EndOfFile) isTerminal()   {}

func NewEmpty() Expr { return Empty{} }
func NewDot() Expr   { return Dot{} }
func NewJunction(exprs ...Expr) Expr {
	if len(exprs) > 1 {
		return Junction{exprs}
	}
	return exprs[0]
}
func NewChoice(exprs ...Expr) Expr {
	if len(exprs) > 1 {
		return Choice{exprs}
	}
	return exprs[0]
}
func NewRepetition(expr Expr) Expr          { return Kleene{expr} }
func NewRepetitionN(expr Expr, n uint) Expr { return Repetition{expr, n} }
func NewNegation(expr Expr) Expr            { return Negation{expr} }
func NewSymbol(s string, attrsOpt ...map[string][]byte) Expr {
	var attrs map[string][]byte
	if len(attrsOpt) > 0 {
		attrs = attrsOpt[0]
	}
	return Symbol{Name: s, Attributes: attrs}
}
func NewOptional(expr Expr) Expr     { return Optional{expr} }
func NewEnsure(expr Expr) Expr       { return Ensure{expr} }
func NewTextToken(token string) Expr { return TextToken{Text: []byte(token)} }
func NewTextPattern(regex string) Expr {
	regex = "^" + strings.TrimPrefix(regex, "^")
	return TextPattern{Expr: regex, Regex: regexp.MustCompile(regex)}
}
func NewAtomPattern(matcher map[string]TextTerminals) Expr {
	return AtomPattern{Matcher: matcher}
}
func NewTokenAttributeMatcher(token string) TextTerminals {
	return TextToken{Text: []byte(token)}
}
func NewPatternAttributeMatcher(regex string) TextTerminals {
	regex = "^" + strings.TrimSuffix(strings.TrimPrefix(regex, "^"), "$") + "$"
	return TextPattern{Expr: regex, Regex: regexp.MustCompile(regex)}
}
