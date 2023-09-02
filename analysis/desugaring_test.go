package analysis

import (
	"github.com/stretchr/testify/assert"
	"gopeg/definition"
	"testing"
)

func TestDesugar(t *testing.T) {
	r := definition.NewRule("S", definition.NewJunction(
		definition.NewOptional(definition.NewSymbol("A")),
		definition.NewEnsure(definition.NewSymbol("B")),
		definition.NewRepetitionN(definition.NewSymbol("C"), 2),
	))
	d := DesugarRule(r)
	assert.NotNil(t, CheckDesugaredRule(r))
	assert.Nil(t, CheckDesugaredRule(d))
	assert.Equal(t, d,
		definition.NewRule("S", definition.Junction{Exprs: []definition.Expr{
			definition.Choice{Exprs: []definition.Expr{
				definition.Symbol{Name: "A"},
				definition.Empty{},
			}},
			definition.Negation{Expr: definition.Negation{Expr: definition.Symbol{Name: "B"}}},
			definition.Junction{
				Exprs: []definition.Expr{
					definition.Symbol{Name: "C"},
					definition.Symbol{Name: "C"},
					definition.Kleene{Expr: definition.Symbol{Name: "C"}},
				},
			},
		}}),
	)
	t.Logf("\ninitial:\n%v\ndesugared:\n%v", r, d)
}
