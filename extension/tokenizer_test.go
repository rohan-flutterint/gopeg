package extension

import (
	"github.com/stretchr/testify/assert"
	"gopeg/parser"
	"testing"
)

func TestControls(t *testing.T) {
	text := `Text: @sof "a" @eof`
	peg, err := parser.ParseText(PegTokenizerRules, PegText, []byte(text))
	assert.Equal(t, len(text), peg.Segment.Length())
	assert.Nil(t, err)
	attrs := make([]map[string]string, 0)
	for _, children := range peg.Children {
		attrs = append(attrs, map[string]string{children.Atom.Symbol: string(children.Atom.SelectText())})
	}
	assert.Equal(t, attrs, []map[string]string{
		{"Token": "Text"},
		{"Control": ":"},
		{"BuiltinSymbol": "@sof"},
		{"String": `"a"`},
		{"BuiltinSymbol": "@eof"},
		{"EndOfLine": ``},
	})
}

func TestEscaping(t *testing.T) {
	text := `Text: "\"\'\\\\\\"`
	peg, err := parser.ParseText(PegTokenizerRules, PegText, []byte(text))
	assert.Equal(t, len(text), peg.Segment.Length())
	assert.Nil(t, err)
	attrs := make([]map[string]string, 0)
	for _, children := range peg.Children {
		attrs = append(attrs, map[string]string{children.Atom.Symbol: string(children.Atom.SelectText())})
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
	peg, err := parser.ParseText(PegTokenizerRules, PegText, []byte(text))
	assert.Equal(t, len(text), peg.Segment.Length())
	assert.Nil(t, err)
	attrs := make([]map[string]string, 0)
	for _, children := range peg.Children {
		attrs = append(attrs, map[string]string{children.Atom.Symbol: string(children.Atom.SelectText())})
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
	peg, err := parser.ParseText(PegTokenizerRules, PegText, []byte(text))
	assert.Equal(t, len(text), peg.Segment.Length())
	assert.Nil(t, err)
	attrs := make([]map[string]string, 0)
	for _, children := range peg.Children {
		attrs = append(attrs, map[string]string{children.Atom.Symbol: string(children.Atom.SelectText())})
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
		{"Dot": "."},
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
