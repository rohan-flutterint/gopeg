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
					NewChoice(NewTerminals("+"), NewTerminals("-")),
					NewNonterminal("Product"),
				),
			),
		)),
		NewRule("Product", NewJunction(
			NewNonterminal("Value"),
			NewRepetition(
				NewJunction(
					NewChoice(NewTerminals("*"), NewTerminals("/")),
					NewNonterminal("Value"),
				),
			),
		)),
		NewRule("Digit", NewChoice(
			NewTerminals("0"),
			NewTerminals("1"),
			NewTerminals("2"),
			NewTerminals("3"),
			NewTerminals("4"),
			NewTerminals("5"),
			NewTerminals("6"),
			NewTerminals("7"),
			NewTerminals("8"),
			NewTerminals("9"),
		)),
		NewRule("Value", NewChoice(
			NewJunction(
				NewNonterminal("Digit"),
				NewRepetition(NewNonterminal("Digit")),
			),
			NewJunction(
				NewTerminals("("),
				NewNonterminal("Expr"),
				NewTerminals(")"),
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
