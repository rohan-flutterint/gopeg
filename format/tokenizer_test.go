package format

import (
	"github.com/stretchr/testify/assert"
	"gopeg/parser"
	"testing"
)

func TestEscaping(t *testing.T) {
	text := `Text: "\"\'\\\\\\"`
	peg, err := parser.Parse[byte](PegTokenizerRules, PegText, []byte(text))
	assert.Equal(t, len(text), peg.Length())
	assert.Nil(t, err)
	attrs := make([]map[string]string, 0)
	for _, children := range peg.Children() {
		attrs = append(attrs, children.Attrs().ToMap())
	}
	assert.Equal(t, attrs, []map[string]string{
		{"Token": "Text"},
		{"Control": ":"},
		{"String": `"\"\'\\\\\\"`},
		{"EndOfLine": ``},
	})
}

func TestParentheses(t *testing.T) {
	text := `Text: (
	"1" / 
	"2" / 
	"3"
)`
	peg, err := parser.Parse[byte](PegTokenizerRules, PegText, []byte(text))
	assert.Equal(t, len(text), peg.Length())
	assert.Nil(t, err)
	attrs := make([]map[string]string, 0)
	for _, children := range peg.Children() {
		attrs = append(attrs, children.Attrs().ToMap())
	}
	assert.Equal(t, attrs, []map[string]string{
		{"Token": `Text`},
		{"Control": `:`},
		{"Open": `(`},
		{"String": `"1"`},
		{"Control": `/`},
		{"String": `"2"`},
		{"Control": `/`},
		{"String": `"3"`},
		{"Close": `)`},
		{"EndOfLine": ``},
	})
}

func TestPegTokenizer(t *testing.T) {
	text := `Text: (=~"[\n\t\r ]+" / "=~" Regex:String / String / =~"//[^\n]+" / Token / Control:.)*
String: =~"'(\\.|[^']*)'" /*
	this is multiline comment
*/
String: =~'"(\\.|[^\"]*)"' // single line comment
Token: =~"[a-zA-Z][0-9a-zA-Z_]+"`
	peg, err := parser.Parse[byte](PegTokenizerRules, PegText, []byte(text))
	assert.Equal(t, len(text), peg.Length())
	assert.Nil(t, err)
	attrs := make([]map[string]string, 0)
	for _, children := range peg.Children() {
		attrs = append(attrs, children.Attrs().ToMap())
	}
	assert.Equal(t, attrs, []map[string]string{
		{"Token": "Text"},
		{"Control": ":"},
		{"Open": "("},
		{"Regex": `"[\n\t\r ]+"`},
		{"Control": "/"},
		{"String": `"=~"`},
		{"Token": "Regex"},
		{"Control": ":"},
		{"Token": "String"},
		{"Control": "/"},
		{"Token": "String"},
		{"Control": "/"},
		{"Regex": `"//[^\n]+"`},
		{"Control": "/"},
		{"Token": "Token"},
		{"Control": "/"},
		{"Token": "Control"},
		{"Control": ":"},
		{"Any": "."},
		{"Close": ")"},
		{"Control": "*"},
		{"EndOfLine": "\n"},

		{"Token": "String"},
		{"Control": ":"},
		{"Regex": `"'(\\.|[^']*)'"`},
		{"EndOfLine": "\n"},

		{"Token": "String"},
		{"Control": ":"},
		{"Regex": `'"(\\.|[^\"]*)"'`},
		{"EndOfLine": "\n"},

		{"Token": "Token"},
		{"Control": ":"},
		{"Regex": `"[a-zA-Z][0-9a-zA-Z_]+"`},
		{"EndOfLine": ""},
	})
}
