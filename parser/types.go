package parser

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	Expr        interface{}
	terminals   struct{ Terminals string }
	nonterminal struct{ NonTerminal string }
	junction    struct{ Exprs []Expr }
	choice      struct{ Exprs []Expr }
	repetition  struct{ Expr Expr }
	negation    struct{ Expr Expr }

	optional struct{ Expr Expr }
	ensure   struct{ Expr Expr }
)

type Rule struct {
	Name string
	Expr Expr
}

type ParsingNode interface {
	NonTerminal() string
	Range() (start, end int)
	Children() []ParsingNode
}

type parsingNode struct {
	nonTerminal string
	start, end  int
	children    []ParsingNode
}

func (p *parsingNode) NonTerminal() string     { return p.nonTerminal }
func (p *parsingNode) Range() (start, end int) { return p.start, p.end }
func (p *parsingNode) Children() []ParsingNode { return p.children }

type Rules []Rule

func NewRule(name string, expression Expr) Rule { return Rule{name, expression} }

func NewTerminals(ts string) Expr    { return terminals{ts} }
func NewJunction(exprs ...Expr) Expr { return junction{exprs} }
func NewChoice(exprs ...Expr) Expr   { return choice{exprs} }
func NewRepetition(expr Expr) Expr   { return repetition{expr} }
func NewNegation(expr Expr) Expr     { return negation{expr} }
func NewNonterminal(nt string) Expr  { return nonterminal{nt} }
func NewOptional(expr Expr) Expr     { return optional{expr} }
func NewEnsure(expr Expr) Expr       { return ensure{expr} }

func NonTerminalName(expr Expr) (string, bool) {
	switch peg := expr.(type) {
	case nonterminal:
		return peg.NonTerminal, true
	default:
		return "", false
	}
}

func StringExpression(expr Expr) string {
	switch peg := expr.(type) {
	case terminals:
		return fmt.Sprintf("'%v'", peg.Terminals)
	case nonterminal:
		return fmt.Sprintf("%v", peg.NonTerminal)
	case repetition:
		return fmt.Sprintf("(%v)*", StringExpression(peg.Expr))
	case negation:
		return fmt.Sprintf("!(%v)", StringExpression(peg.Expr))
	case optional:
		return fmt.Sprintf("(%v)?", StringExpression(peg.Expr))
	case ensure:
		return fmt.Sprintf("&(%v)", StringExpression(peg.Expr))
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
