package definition

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRuleString(t *testing.T) {
	rules := Rules{
		NewRule("Digit", NewTextPattern("[0-9]")),
		NewRule("Letter", NewTextPattern("[a-z]")),
	}
	require.Equal(t, `Digit: =~"^[0-9]"
Letter: =~"^[a-z]"
`, rules.String())
}

func TestRuleCombination(t *testing.T) {
	a := Rules{NewRule("Digit", NewTextPattern("[0-9]"))}
	b := Rules{NewRule("Letter", NewTextPattern("[a-z]"))}
	c := a.Combine(b)
	require.Len(t, a, 1)
	require.Len(t, b, 1)
	require.Len(t, c, 2)
	require.Equal(t, "Digit", c[0].Name)
	require.Equal(t, "Letter", c[1].Name)
}
