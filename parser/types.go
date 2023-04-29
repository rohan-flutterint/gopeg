package parser

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type (
	Expr interface{}

	symbol   struct{ Symbol string }
	junction struct{ Exprs []Expr }
	choice   struct{ Exprs []Expr }
	kleene   struct{ Expr Expr }
	negation struct{ Expr Expr }

	optional   struct{ Expr Expr }
	ensure     struct{ Expr Expr }
	repetition struct {
		Expr Expr
		min  int
	}
)

type Rule struct {
	Name string
	Expr Expr
}

type Attrs map[string]any

func (a Attrs) ToMap() map[string]string {
	m := make(map[string]string, len(a))
	for key, value := range a {
		m[key] = string(value.([]byte))
	}
	return m
}

type ParsingNode interface {
	Attrs() Attrs
	Symbol() string
	Range() (start, end int)
	Length() int
	Children() []ParsingNode
}

func Traverse(n ParsingNode, enter, exit func(node ParsingNode)) {
	if enter != nil {
		enter(n)
	}
	for _, c := range n.Children() {
		Traverse(c, enter, exit)
	}
	if exit != nil {
		exit(n)
	}
}

type Parsing struct {
	Tokens []string
	Meta   []PegBindRef
}

func Extract(n ParsingNode, extract []string) Parsing {
	parsing := Parsing{
		Tokens: make([]string, 0),
		Meta:   make([]PegBindRef, 0),
	}
	extractMap := make(map[string]struct{}, 0)
	for _, e := range extract {
		extractMap[e] = struct{}{}
	}
	Traverse(n,
		func(node ParsingNode) {
			symbol := node.Symbol()
			if _, ok := extractMap[symbol]; ok {
				start, end := node.Range()
				parsing.Tokens = append(parsing.Tokens, symbol)
				parsing.Meta = append(parsing.Meta, PegBindRef{start, end})
			}
		},
		nil,
	)
	return parsing
}

type parsingNode struct {
	symbol     string
	start, end int
	attrs      Attrs
	children   []ParsingNode
}

func NewParsingNode[T any](symbol string, start, end int, text []T) *parsingNode {
	return &parsingNode{
		symbol:   symbol,
		start:    start,
		end:      end,
		children: nil,
		attrs:    Attrs{symbol: text[start:end]},
	}
}

func (p *parsingNode) Symbol() string          { return p.symbol }
func (p *parsingNode) Range() (start, end int) { return p.start, p.end }
func (p *parsingNode) Length() int             { return p.end - p.start }
func (p *parsingNode) Children() []ParsingNode { return p.children }
func (p *parsingNode) Attrs() Attrs            { return p.attrs }

type Rules []Rule

func (rs Rules) Combine(bs ...Rules) Rules {
	for _, b := range bs {
		rs = append(rs, b...)
	}
	return rs
}

func NewRule(name string, expression Expr) Rule { return Rule{name, expression} }

func NewEmpty() Expr { return empty{} }
func NewAny() Expr   { return dot{} }

func NewJunction(exprs ...Expr) Expr {
	if len(exprs) > 1 {
		return junction{exprs}
	}
	return exprs[0]
}
func NewChoice(exprs ...Expr) Expr {
	if len(exprs) > 1 {
		return choice{exprs}
	}
	return exprs[0]
}
func NewRepetition(expr Expr) Expr         { return kleene{expr} }
func NewRepetitionN(expr Expr, n int) Expr { return repetition{expr, n} }
func NewNegation(expr Expr) Expr           { return negation{expr} }
func NewSymbol(s string) Expr              { return symbol{s} }
func NewOptional(expr Expr) Expr           { return optional{expr} }
func NewEnsure(expr Expr) Expr             { return ensure{expr} }
func NewToken(token string) Expr           { return byteSequence{value: []byte(token)} }
func NewMatch(regex string) Expr {
	return byteMatcher{expr: regex, regex: regexp.MustCompile("^" + strings.TrimPrefix(regex, "^"))}
}
func NewAttr(as ...Attrs) Expr {
	matchers := make([]attrMatcher, 0)
	for _, a := range as {
		matcher := make(attrMatcher)
		for key, value := range a {
			if value == nil {
				matcher[key] = nil
			} else {
				matcher[key] = value.(terminals)
			}
		}
		matchers = append(matchers, matcher)
	}
	return attrSequenceMatcher(matchers)
}

func SymbolName(expr Expr) (string, bool) {
	switch peg := expr.(type) {
	case symbol:
		return peg.Symbol, true
	default:
		return "", false
	}
}

func StringExpression(expr Expr) string {
	switch peg := expr.(type) {
	case terminals:
		return fmt.Sprintf("%v", peg.String())
	case symbol:
		return fmt.Sprintf("%v", peg.Symbol)
	case kleene:
		return fmt.Sprintf("(%v)*", StringExpression(peg.Expr))
	case negation:
		return fmt.Sprintf("!(%v)", StringExpression(peg.Expr))
	case optional:
		return fmt.Sprintf("(%v)?", StringExpression(peg.Expr))
	case ensure:
		return fmt.Sprintf("&(%v)", StringExpression(peg.Expr))
	case repetition:
		return fmt.Sprintf("(%v){%v,}", StringExpression(peg.Expr), peg.min)
	case junction:
		names := make([]string, 0, len(peg.Exprs))
		for _, expr := range peg.Exprs {
			names = append(names, StringExpression(expr))
		}
		return fmt.Sprintf("(%v)", strings.Join(names, " "))
	case choice:
		names := make([]string, 0, len(peg.Exprs))
		for _, expr := range peg.Exprs {
			names = append(names, StringExpression(expr))
		}
		return fmt.Sprintf("(%v)", strings.Join(names, " / "))
	default:
		panic(fmt.Errorf("unexpected peg expression type: %v", reflect.TypeOf(expr)))
	}
}

func (r Rule) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%v <- ", r.Name))
	b.WriteString(StringExpression(r.Expr))
	return b.String()
}

func (rs Rules) String() string {
	b := strings.Builder{}
	for _, r := range rs {
		b.WriteString(fmt.Sprintf("%v\n", r))
	}
	return b.String()
}
