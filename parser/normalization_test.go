package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleNormalization(t *testing.T) {
	r := Rules{NewRule("r", NewChoice(
		NewJunction(NewSymbol("A"), NewNegation(NewChoice(NewSymbol("B"), NewSymbol("C")))),
		NewJunction(NewSymbol("C"), NewSymbol("D")),
	))}
	n, _ := r.normalize()
	assert.Equal(t, r.normalized(), false)
	assert.Equal(t, n.normalized(), true)
	assert.Equal(t, n, Rules{
		Rule{"r#0", NewChoice(NewSymbol("r#1"), NewSymbol("r#4"))},
		Rule{"r#1", NewJunction(NewSymbol("A#0"), NewSymbol("r#2"))},
		Rule{"r#2", NewNegation(NewSymbol("r#3"))},
		Rule{"r#3", NewChoice(NewSymbol("B#0"), NewSymbol("C#0"))},
		Rule{"r#4", NewJunction(NewSymbol("C#0"), NewSymbol("D#0"))},
	})
	t.Logf("\ninitial:\n%vdeconstructed:\n%v", r, n)
}

func TestNormalizationWithTerminals(t *testing.T) {
	r := Rules{NewRule("r", NewChoice(
		NewJunction(NewSymbol("A"), NewNegation(NewChoice(NewSymbol("B"), NewSymbol("C")))),
		NewJunction(NewToken("hi"), NewSymbol("D")),
	))}
	n, _ := r.normalize()
	assert.Equal(t, r.normalized(), false)
	assert.Equal(t, n.normalized(), true)
	assert.Equal(t, n, Rules{
		Rule{"r#0", NewChoice(NewSymbol("r#1"), NewSymbol("r#4"))},
		Rule{"r#1", NewJunction(NewSymbol("A#0"), NewSymbol("r#2"))},
		Rule{"r#2", NewNegation(NewSymbol("r#3"))},
		Rule{"r#3", NewChoice(NewSymbol("B#0"), NewSymbol("C#0"))},
		Rule{"r#4", NewJunction(NewToken("hi"), NewSymbol("D#0"))},
	})
	t.Logf("\ninitial:\n%vdeconstructed:\n%v", r, n)
}
