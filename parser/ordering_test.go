package parser

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestOrdering(t *testing.T) {
	rs := Rules{
		NewRule("Expr", NewSymbol("Sum")),
		NewRule("Sum", NewJunction(
			NewSymbol("Product"),
			NewRepetition(
				NewJunction(
					NewChoice(NewToken("+"), NewToken("-")),
					NewSymbol("Product"),
				),
			),
		)),
		NewRule("Product", NewJunction(
			NewSymbol("Value"),
			NewRepetition(
				NewJunction(
					NewChoice(NewToken("*"), NewToken("/")),
					NewSymbol("Value"),
				),
			),
		)),
		NewRule("Digit", NewChoice(
			NewToken("0"),
			NewToken("1"),
			NewToken("2"),
			NewToken("3"),
			NewToken("4"),
			NewToken("5"),
			NewToken("6"),
			NewToken("7"),
			NewToken("8"),
			NewToken("9"),
		)),
		NewRule("Value", NewChoice(
			NewJunction(
				NewSymbol("Digit"),
				NewRepetition(NewSymbol("Digit")),
			),
			NewJunction(
				NewToken("("),
				NewSymbol("Expr"),
				NewToken(")"),
			),
		)),
	}
	n, m := rs.normalize()
	t.Logf("\nrules:\n%v", n)
	order, p, err := n.order()
	t.Logf("order: %v, err: %v", order, err)
	assert.Equal(t, sort.IsSorted(sort.IntSlice([]int{
		p[m.Forward.get("Expr")],
		p[m.Forward.get("Sum")],
		p[m.Forward.get("Product")],
		p[m.Forward.get("Value")],
		p[m.Forward.get("Digit")],
	})), true)
}

func TestCycle(t *testing.T) {
	rs := Rules{
		NewRule("A", NewJunction(NewSymbol("B"), NewSymbol("A"))),
		NewRule("B", NewToken("hello")),
	}
	_, _, err := rs.order()
	assert.Nil(t, err)
}
