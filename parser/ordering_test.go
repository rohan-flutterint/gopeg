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
	//				NewChoice(NewTerminals("+"), NewTerminals("-")),
	//				NewNonterminal("Product"),
	//			),
	//		),
	//	)),
	//	NewRule("Product", NewJunction(
	//		NewNonterminal("Value"),
	//		NewRepetition(
	//			NewJunction(
	//				NewChoice(NewTerminals("*"), NewTerminals("/")),
	//				NewNonterminal("Value"),
	//			),
	//		),
	//	)),
	//	NewRule("Digit", NewChoice(
	//		NewTerminals("0"),
	//		NewTerminals("1"),
	//		NewTerminals("2"),
	//		NewTerminals("3"),
	//		NewTerminals("4"),
	//		NewTerminals("5"),
	//		NewTerminals("6"),
	//		NewTerminals("7"),
	//		NewTerminals("8"),
	//		NewTerminals("9"),
	//	)),
	//	NewRule("Value", NewChoice(
	//		NewJunction(
	//			NewNonterminal("Digit"),
	//			NewRepetition(NewNonterminal("Digit")),
	//		),
	//		NewJunction(
	//			NewTerminals("("),
	//			NewNonterminal("Expr"),
	//			NewTerminals(")"),
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
		NewRule("B", NewTerminals("hello")),
	}
	_, _, err := rs.order()
	assert.Nil(t, err)
}
