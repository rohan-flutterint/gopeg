package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArithmetic(t *testing.T) {
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
		NewRule("Digit", NewMatch("[0-9]")),
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
	text := "10+2"
	n1, err := Parse(rs, "Expr", []byte(text))
	assert.Nil(t, err)
	start, end := n1.Range()
	assert.Equal(t, end-start, 4)
	t.Logf("\n%v\n", StringParsingNode(n1, []byte(text)))

	text = "1+(2*4-5)-10*44"
	n2, err := Parse(rs, "Expr", []byte(text))
	assert.Nil(t, err)
	start, end = n2.Range()
	assert.Equal(t, end-start, 15)
	t.Logf("\n%v\n", StringParsingNode(n2, []byte(text)))
}
