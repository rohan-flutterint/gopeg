package format

import (
	"github.com/stretchr/testify/assert"
	"gopeg/parser"
	"testing"
)

func TestName(t *testing.T) {
	text := `B`
	attrs := make([]map[string]any, 0)
	{
		peg, err := parser.Parse[byte](PegTokenizerRules, PegText, []byte(text))
		assert.Nil(t, err)
		assert.Equal(t, len(text), peg.Length())
		for _, children := range peg.Children() {
			attrs = append(attrs, children.Attrs())
		}
		t.Logf("peg:\n%v", parser.StringParsingNode(peg, []byte(text)))
	}

	{
		peg, err := parser.Parse[map[string]any](PegGrammarRules, PegRule, attrs)
		assert.Nil(t, err)
		t.Logf("peg: %v", parser.StringParsingNode(peg, attrs))
	}
}
