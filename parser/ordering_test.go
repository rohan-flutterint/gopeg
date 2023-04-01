package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrdering(t *testing.T) {
	//rs := Rules{
	//	NewRule("Expr", NewNonterminal("Sum")),
	//	NewRule("Sum", NewJunction(
	//		NewNonterminal("Product"),
	//		NewRepetition(
	//			NewJunction(
	//				NewChoice(NewToken("+"), NewToken("-")),
	//				NewNonterminal("Product"),
	//			),
	//		),
	//	)),
	//	NewRule("Product", NewJunction(
	//		NewNonterminal("Value"),
	//		NewRepetition(
	//			NewJunction(
	//				NewChoice(NewToken("*"), NewToken("/")),
	//				NewNonterminal("Value"),
	//			),
	//		),
	//	)),
	//	NewRule("Digit", NewChoice(
	//		NewToken("0"),
	//		NewToken("1"),
	//		NewToken("2"),
	//		NewToken("3"),
	//		NewToken("4"),
	//		NewToken("5"),
	//		NewToken("6"),
	//		NewToken("7"),
	//		NewToken("8"),
	//		NewToken("9"),
	//	)),
	//	NewRule("Value", NewChoice(
	//		NewJunction(
	//			NewNonterminal("Digit"),
	//			NewRepetition(NewNonterminal("Digit")),
	//		),
	//		NewJunction(
	//			NewToken("("),
	//			NewNonterminal("Expr"),
	//			NewToken(")"),
	//		),
	//	)),
	//}
	//n, m := rs.normalize()
	//t.Logf("\nrules:\n%v", n)
	//order, p, err := n.order()
	//t.Logf("order: %v, err: %v", order, err)
	//assert.Equal(t, sort.IsSorted(sort.IntSlice([]int{p[m["Expr"]], p[m["Sum"]], p[m["Product"]], p[m["Value"]], p[m["Digit"]]})), true)
}

func TestCycle(t *testing.T) {
	rs := Rules{
		NewRule("A", NewJunction(NewNonterminal("B"), NewNonterminal("A"))),
		NewRule("B", NewToken("hello")),
	}
	_, _, err := rs.order()
	assert.Nil(t, err)
}
