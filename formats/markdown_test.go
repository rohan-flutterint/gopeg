package formats

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

//go:embed commonmark-0.30.json
var markdownExamples []byte

type MarkdownExample struct {
	Markdown  string `json:"markdown"`
	Html      string `json:"html"`
	Example   int    `json:"example"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	Section   string `json:"section"`
}

func TestMarkdown(t *testing.T) {
	var examples []MarkdownExample
	require.Nil(t, json.Unmarshal(markdownExamples, &examples))
	sections := make(map[string][]MarkdownExample)
	sectionNames := make([]string, 0)
	for _, example := range examples {
		if _, ok := sections[example.Section]; !ok {
			sections[example.Section] = make([]MarkdownExample, 0)
			sectionNames = append(sectionNames, example.Section)
		}
		sections[example.Section] = append(sections[example.Section], example)
	}
	for _, sectionName := range sectionNames {
		t.Run(sectionName, func(t *testing.T) {
			for _, example := range sections[sectionName] {
				t.Run(fmt.Sprintf("%v(%v)", example.Example, example.Markdown), func(t *testing.T) {
					html, err := Markdown2Html(example.Markdown)
					require.Nil(t, err)
					require.Equal(t, example.Html, html)
				})
			}
		})
	}
}
