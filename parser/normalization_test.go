package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleNormalization(t *testing.T) {
	r := Rules{NewRule("r", NewChoice(
		NewJunction(NewNonterminal("A"), NewNegation(NewChoice(NewNonterminal("B"), NewNonterminal("C")))),
		NewJunction(NewNonterminal("C"), NewNonterminal("D")),
	))}
	n, _ := r.normalize()
	assert.Equal(t, r.normalized(), false)
	assert.Equal(t, n.normalized(), true)
	assert.Equal(t, n, Rules{
		Rule{"r#0", NewChoice(NewNonterminal("r#1"), NewNonterminal("r#4"))},
		Rule{"r#1", NewJunction(NewNonterminal("A#0"), NewNonterminal("r#2"))},
		Rule{"r#2", NewNegation(NewNonterminal("r#3"))},
		Rule{"r#3", NewChoice(NewNonterminal("B#0"), NewNonterminal("C#0"))},
		Rule{"r#4", NewJunction(NewNonterminal("C#0"), NewNonterminal("D#0"))},
	})
	t.Logf("\ninitial:\n%vdeconstructed:\n%v", r, n)
}

func TestNormalizationWithTerminals(t *testing.T) {
	r := Rules{NewRule("r", NewChoice(
		NewJunction(NewNonterminal("A"), NewNegation(NewChoice(NewNonterminal("B"), NewNonterminal("C")))),
		NewJunction(NewToken("hi"), NewNonterminal("D")),
	))}
	n, _ := r.normalize()
	assert.Equal(t, r.normalized(), false)
	assert.Equal(t, n.normalized(), true)
	assert.Equal(t, n, Rules{
		Rule{"r#0", NewChoice(NewNonterminal("r#1"), NewNonterminal("r#4"))},
		Rule{"r#1", NewJunction(NewNonterminal("A#0"), NewNonterminal("r#2"))},
		Rule{"r#2", NewNegation(NewNonterminal("r#3"))},
		Rule{"r#3", NewChoice(NewNonterminal("B#0"), NewNonterminal("C#0"))},
		Rule{"r#4", NewJunction(NewToken("hi"), NewNonterminal("D#0"))},
	})
	t.Logf("\ninitial:\n%vdeconstructed:\n%v", r, n)
}
