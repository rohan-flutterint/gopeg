package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {
	rs := Rules{
		NewRule("Expr", NewNonterminal("Sum")),
		NewRule("Sum", NewJunction(
			NewNonterminal("Product"),
			NewRepetition(
				NewJunction(
					NewChoice(NewToken("+"), NewToken("-")),
					NewNonterminal("Product"),
				),
			),
		)),
		NewRule("Product", NewJunction(
			NewNonterminal("Value"),
			NewRepetition(
				NewJunction(
					NewChoice(NewToken("*"), NewToken("/")),
					NewNonterminal("Value"),
				),
			),
		)),
		NewRule("Digit", NewInterval('0', '9')),
		NewRule("Value", NewChoice(
			NewJunction(
				NewNonterminal("Digit"),
				NewRepetition(NewNonterminal("Digit")),
			),
			NewJunction(
				NewToken("("),
				NewNonterminal("Expr"),
				NewToken(")"),
			),
		)),
	}
	text := "10+2"
	n1, err := rs.parse("Expr", text)
	assert.Nil(t, err)
	start, end := n1.Range()
	assert.Equal(t, end-start, 4)
	t.Logf("\n%v\n", StringParsingNode(n1, text))

	text = "1+(2*4-5)-10*44"
	n2, err := rs.parse("Expr", text)
	assert.Nil(t, err)
	start, end = n2.Range()
	assert.Equal(t, end-start, 15)
	t.Logf("\n%v\n", StringParsingNode(n2, text))
}
