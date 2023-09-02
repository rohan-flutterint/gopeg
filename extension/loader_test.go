package extension

import (
	"github.com/stretchr/testify/assert"
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
