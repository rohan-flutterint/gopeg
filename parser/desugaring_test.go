package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDesugar(t *testing.T) {
	r := NewRule("S", NewJunction(
		NewOptional(NewSymbol("A")),
		NewEnsure(NewSymbol("B")),
		NewRepetition(NewSymbol("C")),
	))
	d := r.desugar()
	assert.Equal(t, d,
		NewRule("S", NewJunction(
			NewChoice(NewSymbol("A"), NewEmpty()),
			NewNegation(NewNegation(NewSymbol("B"))),
			NewRepetition(NewSymbol("C")),
		)),
	)
	t.Logf("\ninitial:\n%v\ndesugared:\n%v", r, d)
}
