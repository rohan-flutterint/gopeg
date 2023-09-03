package formats

import (
	_ "embed"
	"fmt"
	"gopeg/definition"
	"gopeg/extension"
	"gopeg/parser"
	"regexp"
	"strings"
)

//go:embed djot.peg
var djotPegGrammar string

var sectionIdIgnoreRegex = regexp.MustCompile(`[^\w\s]+`)
var sectionIdReplaceRegex = regexp.MustCompile(`\s+`)

func DjotSectionId(node *parser.ParsingNode) string {
	var builder strings.Builder
	node.Traverse(func(node *parser.ParsingNode, next func(nodes []*parser.ParsingNode)) {
		if node.Atom.Symbol == "InlineText" && len(node.Children) == 0 {
			builder.WriteString(node.Atom.SelectString())
		} else {
			next(node.Children)
		}
	})
	id := builder.String()
	id = sectionIdIgnoreRegex.ReplaceAllString(id, "")
	id = sectionIdReplaceRegex.ReplaceAllString(id, "-")
	return id
}

var djotRules definition.Rules

func init() {
	rules, err := extension.Load(djotPegGrammar)
	if err != nil {
		panic(fmt.Errorf("unable to load djot grammar: %w", err))
	}
	djotRules = rules
}

func Djot2Html(markdown string) (string, error) {
	node, err := parser.ParseText(djotRules, "Djot", []byte(markdown))
	if err != nil {
		return "", fmt.Errorf("unable to parse djot text: %w", err)
	}
	if node.Segment.Length() != len(markdown) {
		return "", fmt.Errorf("unable to fully parse markdown (valid first %v bytes)", node.Segment.Length())
	}
	var html strings.Builder
	footnotes := make(map[string]string)
	node.Traverse(func(node *parser.ParsingNode, next func(nodes []*parser.ParsingNode)) {
		switch node.Atom.Symbol {
		case "FootnoteReference":
			footnotes[node.MustSelectBySymbol("Name").Atom.SelectString()] = node.MustSelectBySymbol("Link").Atom.SelectString()
		case "Link", "Reference":
		default:
			next(node.Children)
		}
	})
	node.Traverse(func(node *parser.ParsingNode, next func(nodes []*parser.ParsingNode)) {
		switch node.Atom.Symbol {
		case "CodeBlock":
			html.WriteString("<pre>")
			if lang, ok := node.TrySelectBySymbol("Language"); ok {
				html.WriteString(fmt.Sprintf("<code lang=\"%v\">", lang.Atom.SelectString()))
			} else {
				html.WriteString("<code>")
			}
			next(node.Children)
			html.WriteString("</code></pre>\n")
		case "SimpleHtml":
			newline := ""
			if _, ok := node.Atom.Attributes["Block"]; ok {
				newline = "\n"
			}
			if len(node.Children) == 0 {
				html.WriteString(fmt.Sprintf("<%v>%v", string(node.Atom.Attributes["Tag"]), newline))
			} else {
				html.WriteString(fmt.Sprintf("<%v>%v", string(node.Atom.Attributes["Tag"]), newline))
				next(node.Children)
				html.WriteString(fmt.Sprintf("</%v>%v", string(node.Atom.Attributes["Tag"]), newline))
			}
		case "Div":
			html.WriteString(fmt.Sprintf("<div class=\"%v\">\n", node.MustSelectBySymbol("Class").Atom.SelectString()))
			next(node.Children)
			html.WriteString("</div>\n")
		case "Heading":
			headingId := strings.ReplaceAll(DjotSectionId(node.Children[1]), " ", "-")
			level := len(node.MustSelectBySymbol("Level").Atom.SelectText())
			html.WriteString(fmt.Sprintf("<section id=\"%v\">\n", headingId))
			html.WriteString(fmt.Sprintf(`<h%v>`, level))
			next(node.Children[1:2])
			html.WriteString(fmt.Sprintf("</h%v>\n", level))
			next(node.Children[2:])
			html.WriteString("</section>\n")
		case "Paragraph":
			html.WriteString("<p>")
			next(node.Children)
			html.WriteString("</p>\n")
		case "Verbatim":
			code := parser.ConcatNodeTexts(node.Children...)
			if strings.HasPrefix(code, " `") && strings.HasSuffix(code, "` ") {
				code = code[1 : len(code)-1]
			}
			html.WriteString(fmt.Sprintf("<%v", string(node.Atom.Attributes["Tag"])))
			if class, ok := node.Atom.Attributes["Class"]; ok {
				classString := string(class)
				html.WriteString(fmt.Sprintf(" class=\"%v\"", classString))
				if classString == "math inline" {
					code = "\\(" + code + "\\)"
				} else if classString == "math display" {
					code = "\\[" + code + "\\]"
				}
			}
			html.WriteString(">")
			html.WriteString(code)
			html.WriteString(fmt.Sprintf("</%v>", string(node.Atom.Attributes["Tag"])))
		case "Link", "AutoLink":
			_, isImage := node.Atom.Attributes["Image"]
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
		case "InlineText":
			if len(node.Children) > 0 {
				next(node.Children)
			} else {
				html.WriteString(string(node.Atom.SelectText()))
			}
		case "LineBreak":
			html.WriteString("<br>\n")
		case "Symbol":
			value := node.Atom.SelectString()
			if value == "+1" {
				html.WriteString("👍")
			} else if value == "smiley" {
				html.WriteString("😃")
			} else {
				html.WriteString(":" + value + ":")
			}
		case "FootnoteReference":
		default:
			next(node.Children)
		}
	})
	return html.String(), nil
}
