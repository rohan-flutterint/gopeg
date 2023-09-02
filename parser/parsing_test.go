package parser

import (
	"github.com/stretchr/testify/assert"
	"gopeg/definition"
	"testing"
)

func TestArithmetic(t *testing.T) {
	rs := definition.Rules{
		definition.NewRule("Expr", definition.NewSymbol("Sum")),
		definition.NewRule("Sum", definition.NewJunction(
			definition.NewSymbol("Product"),
			definition.NewRepetition(
				definition.NewJunction(
					definition.NewChoice(
						definition.NewTextToken("+"),
						definition.NewTextToken("-"),
					),
					definition.NewSymbol("Product"),
				),
			),
		)),
		definition.NewRule("Product", definition.NewJunction(
			definition.NewSymbol("Value"),
			definition.NewRepetition(
				definition.NewJunction(
					definition.NewChoice(
						definition.NewTextToken("*"),
						definition.NewTextToken("/"),
					),
					definition.NewSymbol("Value"),
				),
			),
		)),
		definition.NewRule("Digit", definition.NewTextPattern("[0-9]")),
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
	text := "10+2"
	n1, err := ParseText(rs, "Expr", []byte(text))
	assert.Nil(t, err)
	t.Logf("%#v\n", n1)
	t.Logf("\n%v\n", StringParsingNode(n1))

	text = "1+(2*4-5)-10*44"
	n2, err := ParseText(rs, "Expr", []byte(text))
	assert.Nil(t, err)
	t.Logf("\n%v\n", StringParsingNode(n2))
}
