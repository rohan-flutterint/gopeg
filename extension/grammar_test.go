package extension

import (
	"github.com/stretchr/testify/assert"
	"gopeg/definition"
	"gopeg/parser"
	"testing"
)

func TestGrammar(t *testing.T) {
	text := `Definitions: Definition*`
	atoms := make([]definition.Atom, 0)
	{
		peg, err := parser.ParseText(PegTokenizerRules, PegText, []byte(text))
		assert.Nil(t, err)
		assert.Equal(t, len(text), peg.Segment.Length())
		for _, children := range peg.Children {
			atoms = append(atoms, children.Atom)
		}
		t.Logf("peg:\n%v", parser.StringParsingNode(peg))
	}

	{
		peg, err := parser.ParseAtoms(PegGrammarRules, PegRule, atoms)
		assert.Nil(t, err)
		t.Logf("peg: %v", parser.StringParsingNode(peg))
	}
}
