package extension

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopeg/parser"
	"os"
	"testing"
)

func TestTokenizer(t *testing.T) {
	tokenizer, err := os.ReadFile("peg-tokenizer.peg")
	assert.Nil(t, err)
	rules, err := Load(string(tokenizer))
	assert.Nil(t, err)
	t.Logf("rules: %v", rules)
}

func TestGrammarLoad(t *testing.T) {
	tokenizer, err := os.ReadFile("peg-grammar.peg")
	assert.Nil(t, err)
	rules, err := Load(string(tokenizer))
	assert.Nil(t, err)
	t.Logf("rules: %v", rules)
}

func TestInlineRules(t *testing.T) {
	rules, err := Load(`A: (Text:=~"[0-9]")+ / (Text:=~"[a-z]")+`)
	require.Nil(t, err)
	{
		node, err := parser.ParseText(rules, "A", []byte(`1234`))
		require.Nil(t, err)
		require.Equal(t, 4, node.Segment.Length())
		require.Len(t, node.Children, 4)
		require.Equal(t, "Text", node.Children[0].Atom.Symbol)
	}
	{
		node, err := parser.ParseText(rules, "A", []byte(`abcd`))
		require.Nil(t, err)
		require.Equal(t, 4, node.Segment.Length())
	}
	{
		node, err := parser.ParseText(rules, "A", []byte(`12ab`))
		require.Nil(t, err)
		require.Equal(t, 2, node.Segment.Length())
	}
}
