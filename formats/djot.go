package formats

import (
	_ "embed"
	"fmt"
	"gopeg/extension"
	"gopeg/parser"
	"strings"
)

//go:embed djot.peg
var djotPegGrammar string

func Djot2Html(markdown string) (string, error) {
	rules, err := extension.Load(djotPegGrammar)
	if err != nil {
		return "", err
	}
	node, err := parser.ParseText(rules, "Djot", []byte(markdown))
	if err != nil {
		return "", err
	}
	if node.Segment.Length() != len(markdown) {
		return "", fmt.Errorf("unable to fully parse markdown (valid first %v bytes)", node.Segment.Length())
	}
	var html strings.Builder
	footnotes := make(map[string]string)
	node.Traverse(func(node *parser.ParsingNode, next func(nodes ...*parser.ParsingNode)) {
		switch node.Atom.Symbol {
		case "Reference":
			footnotes[node.MustSelectBySymbol("Name").Atom.SelectString()] = node.MustSelectBySymbol("Link").Atom.SelectString()
		case "Link":
		default:
			next()
		}
	})
	node.Traverse(func(node *parser.ParsingNode, next func(nodes ...*parser.ParsingNode)) {
		switch node.Atom.Symbol {
		case "Heading":
			level := len(node.MustSelectBySymbol("Level").Atom.SelectText())
			html.WriteString(fmt.Sprintf("<h%v>", level))
			next()
			html.WriteString(fmt.Sprintf("</h%v>\n", level))
		case "Paragraph":
			if string(node.Atom.SelectText()) == "" {
				html.WriteString("\n")
				return
			}
			html.WriteString("<p>")
			next()
			html.WriteString("</p>\n")
		case "EmphasisLight":
			html.WriteString("<em>")
			next()
			html.WriteString("</em>")
		case "EmphasisHeavy":
			html.WriteString("<strong>")
			next()
			html.WriteString("</strong>")
		case "Highlighted":
			html.WriteString("<mark>")
			next()
			html.WriteString("</mark>")
		case "Subscript":
			html.WriteString("<sub>")
			next()
			html.WriteString("</sub>")
		case "Superscript":
			html.WriteString("<sup>")
			next()
			html.WriteString("</sup>")
		case "Insert":
			html.WriteString("<ins>")
			next()
			html.WriteString("</ins>")
		case "Delete":
			html.WriteString("<del>")
			next()
			html.WriteString("</del>")
		case "Verbatim":
			html.WriteString("<code>")
			code := parser.ConcatNodeTexts(node.Children...)
			if strings.HasPrefix(code, " `") && strings.HasSuffix(code, "` ") {
				code = code[1 : len(code)-1]
			}
			html.WriteString(code)
			html.WriteString("</code>")
		case "Link", "AutoLink":
			isImageNode, hasIsImageNode := node.TrySelectBySymbol("IsImage")
			isImage := hasIsImageNode && isImageNode.Atom.SelectString() == "!"
			textNode, hasTextNode := node.TrySelectBySymbol("Text")
			var text string
			if hasTextNode {
				text = parser.ConcatNodeTexts(textNode.FilterBySymbol("InlineText")...)
			}
			urlNode, hasUrl := node.TrySelectBySymbol("Url")
			referenceNode, hasReference := node.TrySelectBySymbol("Reference")
			var src string
			if hasUrl {
				src = parser.ConcatNodeTexts(urlNode.FilterBySymbol("InlineText")...)
			} else if hasReference {
				reference := parser.ConcatNodeTexts(referenceNode)
				if hasReference && reference == "" {
					reference = text
				}
				src = footnotes[reference]
			}
			if text == "" {
				text = src
			}
			if strings.Contains(src, "@") {
				src = "mailto:" + src
			}
			if isImage {
				html.WriteString(fmt.Sprintf(`<img alt="%v" src="%v">`, text, src))
			} else {
				if src != "" {
					html.WriteString(fmt.Sprintf(`<a href="%v">%v</a>`, src, text))
				} else {
					html.WriteString(fmt.Sprintf(`<a>%v</a>`, text))
				}
			}
		case "Pre":
			html.WriteString(fmt.Sprintf("<pre class=\"%v\">%v</pre>\n",
				node.MustSelectBySymbol("Info").Atom.SelectString(),
				node.MustSelectBySymbol("Text").Atom.SelectString(),
			))
		case "InlineText":
			if len(node.Children) > 0 {
				next()
			} else {
				html.WriteString(string(node.Atom.SelectText()))
			}
		case "Reference":
		default:
			next()
		}
	})
	return html.String(), nil
}
