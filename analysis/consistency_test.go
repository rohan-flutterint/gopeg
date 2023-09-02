package analysis

import (
	"github.com/stretchr/testify/require"
	"gopeg/definition"
	"testing"
)

func TestConsistency(t *testing.T) {
	t.Run("only bytes", func(t *testing.T) {
		terminalType, err := CheckRulesConsistency(definition.Rules{
			definition.NewRule("A", definition.NewTextToken("hello")),
			definition.NewRule("B", definition.NewJunction(definition.NewTextToken("world"), definition.NewDot())),
		})
		require.Nil(t, err)
		require.Equal(t, ByteTerminalType, terminalType)
	})
	t.Run("only atoms", func(t *testing.T) {
		terminalType, err := CheckRulesConsistency(definition.Rules{
			definition.NewRule("A", definition.NewAtomPattern(map[string]definition.AttributeMatcher{})),
		})
		require.Nil(t, err)
		require.Equal(t, AtomTerminalType, terminalType)
	})
	t.Run("any terminals type", func(t *testing.T) {
		terminalType, err := CheckRulesConsistency(definition.Rules{
			definition.NewRule("B", definition.NewJunction(definition.NewDot(), definition.NewDot())),
		})
		require.Nil(t, err)
		require.Equal(t, AnyTerminalType, terminalType)
	})
	t.Run("inconsistent expression", func(t *testing.T) {
		_, err := CheckRulesConsistency(definition.Rules{
			definition.NewRule("B", definition.NewJunction(
				definition.NewTextToken("hello"),
				definition.NewAtomPattern(map[string]definition.AttributeMatcher{"A": definition.NewTokenAttributeMatcher("test")}),
			)),
		})
		require.NotNil(t, err)
		t.Log(err)
	})
	t.Run("inconsistent rules", func(t *testing.T) {
		_, err := CheckRulesConsistency(definition.Rules{
			definition.NewRule("A", definition.NewAtomPattern(map[string]definition.AttributeMatcher{})),
			definition.NewRule("A", definition.NewTextToken("hello")),
		})
		require.NotNil(t, err)
		t.Log(err)
	})
}
