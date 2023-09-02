package parser

import (
	"github.com/stretchr/testify/assert"
	"gopeg/analysis"
	"gopeg/definition"
	"sort"
	"testing"
)

func TestOrdering(t *testing.T) {
	rs := definition.Rules{
		definition.NewRule("Expr", definition.NewSymbol("Sum")),
		definition.NewRule("Sum", definition.NewJunction(
			definition.NewSymbol("Product"),
			definition.NewRepetition(
				definition.NewJunction(
					definition.NewChoice(definition.NewTextToken("+"), definition.NewTextToken("-")),
					definition.NewSymbol("Product"),
				),
			),
		)),
		definition.NewRule("Product", definition.NewJunction(
			definition.NewSymbol("Value"),
			definition.NewRepetition(
				definition.NewJunction(
					definition.NewChoice(definition.NewTextToken("*"), definition.NewTextToken("/")),
					definition.NewSymbol("Value"),
				),
			),
		)),
		definition.NewRule("Digit", definition.NewChoice(
			definition.NewTextToken("0"),
			definition.NewTextToken("1"),
			definition.NewTextToken("2"),
			definition.NewTextToken("3"),
			definition.NewTextToken("4"),
			definition.NewTextToken("5"),
			definition.NewTextToken("6"),
			definition.NewTextToken("7"),
			definition.NewTextToken("8"),
			definition.NewTextToken("9"),
		)),
		definition.NewRule("Value", definition.NewChoice(
			definition.NewJunction(
				definition.NewSymbol("Digit"),
				definition.NewRepetition(definition.NewSymbol("Digit")),
			),
			definition.NewJunction(
				definition.NewTextToken("("),
				definition.NewSymbol("Expr"),
				definition.NewTextToken(")"),
			),
		)),
	}
	n, m := analysis.NormalizeRules(rs)
	t.Logf("\nrules:\n%v", n)
	order, p, err := OrderRules(n)
	t.Logf("order: %v, err: %v", order, err)
	assert.True(t, sort.IsSorted(sort.IntSlice([]int{
		p[m.Forward["Expr"]],
		p[m.Forward["Sum"]],
		p[m.Forward["Product"]],
		p[m.Forward["Value"]],
		p[m.Forward["Digit"]],
	})))
}

func TestCycle(t *testing.T) {
	rs := definition.Rules{
		definition.NewRule("A", definition.NewJunction(definition.NewSymbol("B"), definition.NewSymbol("A"))),
		definition.NewRule("B", definition.NewTextToken("hello")),
	}
	_, _, err := OrderRules(rs)
	assert.Nil(t, err)
}
