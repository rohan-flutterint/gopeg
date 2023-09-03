package parser

import (
	"github.com/sivukhin/gopeg/definition"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAttributes(t *testing.T) {
	rs := definition.Rules{
		definition.NewRule("A", definition.NewRepetition(definition.NewChoice(
			definition.NewJunction(definition.NewTextToken("B"), definition.NewSymbol("#B")),
			definition.NewJunction(definition.NewTextToken("C"), definition.NewSymbol("#C")),
		))),
		definition.NewRule("#B", definition.NewSymbol("D", map[string][]byte{"Ctx": []byte("B")})),
		definition.NewRule("#C", definition.NewSymbol("D", map[string][]byte{"Ctx": []byte("C")})),
		definition.NewRule("D", definition.NewTextToken(".")),
	}
	text := "B.B.C.B."
	node, err := ParseText(rs, "A", []byte(text))
	require.Nil(t, err)
	require.Len(t, node.Children, 4)
	require.Equal(t, []byte("B"), node.Children[0].Atom.Attributes["Ctx"])
	require.Equal(t, []byte("B"), node.Children[1].Atom.Attributes["Ctx"])
	require.Equal(t, []byte("C"), node.Children[2].Atom.Attributes["Ctx"])
	require.Equal(t, []byte("B"), node.Children[3].Atom.Attributes["Ctx"])
}

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
