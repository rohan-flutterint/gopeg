package formats

import (
	_ "embed"
	"fmt"
	"gopeg/extension"
	"gopeg/parser"
	"strings"
)

//go:embed markdown-tokenizer.peg
var markdownTokenizer string

func Markdown2Html(markdown string) (string, error) {
	rules, err := extension.Load(markdownTokenizer)
	if err != nil {
		return "", err
	}
	node, err := parser.ParseText(rules, "Markdown", []byte(markdown))
	if err != nil {
		return "", err
	}
	if node.Segment.Length() != len(markdown) {
		return "", fmt.Errorf("unable to fully parse markdown (valid first %v bytes)", node.Segment.Length())
	}
	var html strings.Builder
	for _, child := range node.Children {
		switch child.Atom.Symbol {
		case "Heading":
			level := len(child.MustSelectBySymbol("Level").Atom.SelectText())
			text := string(child.MustSelectBySymbol("HeadingText").Atom.SelectText())
			html.WriteString(fmt.Sprintf("<h%v>%v</h%v>\n", level, text, level))
		case "Paragraph":
			if string(child.Atom.SelectText()) == "" {
				html.WriteString("\n")
				continue
			}
			html.WriteString("<p>")
			child.Traverse(func(node *parser.ParsingNode, next func()) {
				fmt.Printf("%v\n", node.Atom.Symbol)
				switch node.Atom.Symbol {
				case "Emphasis1":
					html.WriteString("<em>")
					next()
					html.WriteString("</em>")
				case "Emphasis2":
					html.WriteString("<strong>")
					next()
					html.WriteString("</strong>")
				case "PlainText":
					html.WriteString(string(node.Atom.SelectText()))
				default:
					next()
				}
			})
			html.WriteString("</p>\n")
		}
	}
	return html.String(), nil
}
