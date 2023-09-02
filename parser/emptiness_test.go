package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopeg/definition"
	"gopeg/analysis"
	"testing"
)

func TestEmptiness(t *testing.T) {
	t.Run("A:A+", func(t *testing.T) {
		rules, _ := analysis.NormalizeRules(analysis.DesugarRules(definition.Rules{
			definition.NewRule("A", definition.NewRepetitionN(definition.NewSymbol("A"), 1)),
		}))
		matchMapValues(t, map[string]bool{"A#0": false}, GetRulesEmptiness(rules))
	})
	t.Run(`A:@empty "a"`, func(t *testing.T) {
		require.Equal(t, map[string]bool{"A": false}, GetRulesEmptiness(definition.Rules{
			definition.NewRule("A", definition.NewJunction(definition.NewEmpty(), definition.NewTextToken("a"))),
		}))
	})
	t.Run(`A:@empty/"a"`, func(t *testing.T) {
		require.Equal(t, map[string]bool{"A": true}, GetRulesEmptiness(definition.Rules{
			definition.NewRule("A", definition.NewChoice(definition.NewEmpty(), definition.NewTextToken("a"))),
		}))
	})
	t.Run(`A:B C|B:"a"+|C:C* "x"`, func(t *testing.T) {
		rules, _ := analysis.NormalizeRules(analysis.DesugarRules(definition.Rules{
			definition.NewRule("A", definition.NewJunction(definition.NewSymbol("B"), definition.NewSymbol("C"))),
			definition.NewRule("B", definition.NewRepetitionN(definition.NewTextToken("a"), 1)),
			definition.NewRule("C", definition.NewJunction(definition.NewRepetition(definition.NewSymbol("C")), definition.NewTextToken("x"))),
		}))
		t.Log(rules)
		matchMapValues(t, map[string]bool{"A#0": false, "B#0": false, "C#0": false}, GetRulesEmptiness(rules))
	})
	t.Run(`A:B C|B:"a"+|C:C*`, func(t *testing.T) {
		rules, _ := analysis.NormalizeRules(analysis.DesugarRules(definition.Rules{
			definition.NewRule("A", definition.NewJunction(definition.NewSymbol("B"), definition.NewSymbol("C"))),
			definition.NewRule("B", definition.NewRepetitionN(definition.NewTextToken("a"), 1)),
			definition.NewRule("C", definition.NewRepetition(definition.NewSymbol("C"))),
		}))
		t.Log(rules)
		matchMapValues(t, map[string]bool{"A#0": false, "B#0": false, "C#0": true}, GetRulesEmptiness(rules))
	})
	t.Run(`A:B C D|B:"a"+|C:"b"+|D:B* C*`, func(t *testing.T) {
		rules, _ := analysis.NormalizeRules(analysis.DesugarRules(definition.Rules{
			definition.NewRule("A", definition.NewJunction(definition.NewSymbol("B"), definition.NewSymbol("C"))),
			definition.NewRule("B", definition.NewRepetitionN(definition.NewTextToken("a"), 1)),
			definition.NewRule("C", definition.NewRepetitionN(definition.NewTextToken("b"), 1)),
			definition.NewRule("D", definition.NewChoice(definition.NewRepetition(definition.NewSymbol("B")), definition.NewRepetition(definition.NewSymbol("C")))),
		}))
		t.Log(rules)
		matchMapValues(t, map[string]bool{"A#0": false, "B#0": false, "C#0": false, "D#0": true}, GetRulesEmptiness(rules))
	})
	t.Run("A:B|B:C|C:D|D:E|E:@empty", func(t *testing.T) {
		rules := definition.Rules{
			definition.NewRule("A", definition.NewSymbol("B")),
			definition.NewRule("B", definition.NewSymbol("C")),
			definition.NewRule("C", definition.NewSymbol("D")),
			definition.NewRule("D", definition.NewSymbol("E")),
			definition.NewRule("E", definition.NewEmpty()),
		}
		t.Log(rules)
		matchMapValues(t, map[string]bool{"A": true, "B": true, "C": true, "D": true, "E": true}, GetRulesEmptiness(rules))
	})
}

func matchMapValues[Key comparable, Value any](t *testing.T, expected, actual map[Key]Value) {
	for key, expectedValue := range expected {
		actualValue, ok := actual[key]
		assert.Truef(t, ok, "key %v is absent", key)
		assert.Equalf(t, expectedValue, actualValue, "key %v: %v != %v", key, expectedValue, actualValue)
	}
}
