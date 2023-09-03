package analysis

import (
	"github.com/stretchr/testify/assert"
	"github.com/sivukhin/gopeg/definition"
	"testing"
)

func TestDuplicateNormalization(t *testing.T) {
	r := definition.Rules{
		definition.NewRule("A", definition.NewTextToken("hello")),
		definition.NewRule("A", definition.NewTextToken("world")),
	}
	n, _ := NormalizeRules(r)
	assert.NotNil(t, CheckNormalizedRules(r))
	assert.Nil(t, CheckNormalizedRules(n))
	assert.Equal(t, n, definition.Rules{
		definition.NewRule("A#0", definition.NewChoice(
			definition.NewTextToken("hello"),
			definition.NewTextToken("world"),
		))},
	)
	t.Log(r, n)
}

func TestSimpleNormalization(t *testing.T) {
	r := definition.Rules{definition.NewRule("r", definition.NewChoice(
		definition.NewJunction(definition.NewSymbol("A"), definition.NewNegation(definition.NewChoice(definition.NewSymbol("B"), definition.NewSymbol("C")))),
		definition.NewJunction(definition.NewSymbol("C"), definition.NewSymbol("D")),
	))}
	n, _ := NormalizeRules(r)
	assert.NotNil(t, CheckNormalizedRules(r))
	assert.Nil(t, CheckNormalizedRules(n))
	assert.Equal(t, n, definition.Rules{
		definition.Rule{Name: "r#0", Expr: definition.NewChoice(definition.NewSymbol("r#3"), definition.NewSymbol("r#4"))},
		definition.Rule{Name: "r#1", Expr: definition.NewChoice(definition.NewSymbol("B#0"), definition.NewSymbol("C#0"))},
		definition.Rule{Name: "r#2", Expr: definition.NewNegation(definition.NewSymbol("r#1"))},
		definition.Rule{Name: "r#3", Expr: definition.NewJunction(definition.NewSymbol("A#0"), definition.NewSymbol("r#2"))},
		definition.Rule{Name: "r#4", Expr: definition.NewJunction(definition.NewSymbol("C#0"), definition.NewSymbol("D#0"))},
	})
	t.Log(r, n)
}

func TestNormalizationWithTerminals(t *testing.T) {
	r := definition.Rules{definition.NewRule("r", definition.NewChoice(
		definition.NewJunction(definition.NewSymbol("A"), definition.NewNegation(definition.NewChoice(definition.NewSymbol("B"), definition.NewSymbol("C")))),
		definition.NewJunction(definition.NewTextToken("hi"), definition.NewSymbol("D")),
	))}
	n, _ := NormalizeRules(r)
	assert.NotNil(t, CheckNormalizedRules(r))
	assert.Nil(t, CheckNormalizedRules(n))
	assert.Equal(t, n, definition.Rules{
		definition.Rule{Name: "r#0", Expr: definition.NewChoice(definition.NewSymbol("r#3"), definition.NewSymbol("r#4"))},
		definition.Rule{Name: "r#1", Expr: definition.NewChoice(definition.NewSymbol("B#0"), definition.NewSymbol("C#0"))},
		definition.Rule{Name: "r#2", Expr: definition.NewNegation(definition.NewSymbol("r#1"))},
		definition.Rule{Name: "r#3", Expr: definition.NewJunction(definition.NewSymbol("A#0"), definition.NewSymbol("r#2"))},
		definition.Rule{Name: "r#4", Expr: definition.NewJunction(definition.NewTextToken("hi"), definition.NewSymbol("D#0"))},
	})
	t.Log(r, n)
}
