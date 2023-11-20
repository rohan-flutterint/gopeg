package highlight

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/sivukhin/gopeg/definition"
	"github.com/sivukhin/gopeg/extension"
	"github.com/sivukhin/gopeg/parser"
)

var (
	//go:embed python-tokenizer.peg
	PythonTokenizer      string
	PythonTokenizerRules definition.Rules
	//go:embed c-tokenizer.peg
	CTokenizer      string
	CTokenizerRules definition.Rules
)

func init() {
	var err error
	PythonTokenizerRules, err = extension.Load(PythonTokenizer)
	if err != nil {
		panic(fmt.Errorf("unable to load PythonTokenizer rules: %w", err))
	}
	CTokenizerRules, err = extension.Load(CTokenizer)
	if err != nil {
		panic(fmt.Errorf("unable to load CTokenizer rules: %w", err))
	}
}

func Highlight(text string, tokenRules definition.Rules) (string, error) {
	tokens, err := parser.ParseText(tokenRules, tokenRules[0].Name, []byte(text))
	if err != nil {
		return "", fmt.Errorf("unable to parse tokens for highlight: root=%v, err=%w", tokenRules[0].Name, err)
	}
	if tokens.Segment.Length() != len(text) {
		return "", fmt.Errorf("unable to highlight whole document, matched only first %v symbols: root=%v, err=%w", tokens.Segment.Length(), tokenRules[0].Name, err)
	}
	result := strings.Builder{}
	tokens.Traverse(func(node *parser.ParsingNode, next func(nodes []*parser.ParsingNode)) {
		if tag := node.Atom.Attributes["tag"]; tag != nil {
			result.WriteString("<")
			result.Write(tag)
			if class := node.Atom.Attributes["class"]; class != nil {
				result.WriteString(" class=\"")
				result.Write(class)
				result.WriteString("\"")
			}
			result.WriteString(">")
		}

		if len(node.Children) > 0 {
			next(node.Children)
		} else {
			result.Write(node.Atom.SelectText())
		}

		if tag := node.Atom.Attributes["tag"]; tag != nil {
			result.WriteString("</")
			result.Write(tag)
			result.WriteString(">")
		}
	})
	return result.String(), nil
}
