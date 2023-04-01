package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDesugar(t *testing.T) {
	r := NewRule("S", NewJunction(
		NewOptional(NewNonterminal("A")),
		NewEnsure(NewNonterminal("B")),
		NewRepetition(NewNonterminal("C")),
	))
	d := r.desugar()
	assert.Equal(t, d,
		NewRule("S", NewJunction(
			NewChoice(NewNonterminal("A"), NewToken("")),
			NewNegation(NewNegation(NewNonterminal("B"))),
			NewRepetition(NewNonterminal("C")),
		)),
	)
	t.Logf("\ninitial:\n%v\ndesugared:\n%v", r, d)
}
