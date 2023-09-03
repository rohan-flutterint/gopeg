package definition

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSymbols(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		symbol, hidden := AnalyzeSymbolName("Text")
		require.Equal(t, "Text", symbol)
		require.False(t, hidden)
	})
	t.Run("hidden", func(t *testing.T) {
		symbol, hidden := AnalyzeSymbolName("#Hidden")
		require.Equal(t, "#Hidden", symbol)
		require.True(t, hidden)
	})
	t.Run("inline", func(t *testing.T) {
		symbol, hidden := AnalyzeSymbolName("Text@2")
		require.Equal(t, "Text", symbol)
		require.False(t, hidden)
	})
	t.Run("hidden-inline", func(t *testing.T) {
		symbol, hidden := AnalyzeSymbolName("#Text@2")
		require.Equal(t, "#Text", symbol)
		require.True(t, hidden)
	})
}
