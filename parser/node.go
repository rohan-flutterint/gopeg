package parser

import (
	"fmt"
	"gopeg/definition"
	"strings"
)

type ParsingNode struct {
	Atom     definition.Atom
	Segment  definition.Segment
	Children []*ParsingNode
}

func (n *ParsingNode) Traverse(process func(node *ParsingNode, next func())) {
	process(n, func() {
		for _, child := range n.Children {
			child.Traverse(process)
		}
	})
}

func (n *ParsingNode) FilterBySymbol(symbol string) []*ParsingNode {
	filtered := make([]*ParsingNode, 0)
	for _, child := range n.Children {
		if child.Atom.Symbol == symbol {
			filtered = append(filtered, child)
		}
	}
	return filtered
}

func (n *ParsingNode) EnsureOnlySymbol(symbol string) []*ParsingNode {
	for _, child := range n.Children {
		if child.Atom.Symbol != symbol {
			panic(fmt.Errorf("unexpected symbol '%v' in node '%v'", symbol, n))
		}
	}
	return n.Children
}

func (n *ParsingNode) EnsureOnlySingle() *ParsingNode {
	if len(n.Children) != 1 {
		panic(fmt.Errorf("unexpected amount of children in node '%v': %v != 1", n, len(n.Children)))
	}
	return n.Children[0]
}

func (n *ParsingNode) TrySelectBySymbol(symbol string) (*ParsingNode, bool) {
	var target *ParsingNode
	count := 0
	for _, child := range n.Children {
		if child.Atom.Symbol == symbol {
			count++
			target = child
		}
	}
	if count > 1 {
		panic(fmt.Errorf("unexpected amount of symbols '%v' in node '%v': %v > 1", symbol, n, count))
	}
	if count == 1 {
		return target, true
	}
	return nil, false
}

func (n *ParsingNode) MustSelectBySymbol(symbol string) *ParsingNode {
	target, ok := n.TrySelectBySymbol(symbol)
	if !ok {
		panic(fmt.Errorf("symbol '%v' not found in node '%v'", symbol, n))
	}
	return target
}

func NewParsingNode[T any](symbol string, attributes map[string][]byte, data []T, segment definition.Segment) ParsingNode {
	text, textOk := any(data).([]byte)
	atoms, atomsOk := any(data).([]definition.Atom)
	var textSelector definition.Segments
	if textOk {
		textSelector = definition.BuildSegments(segment)
	} else if atomsOk {
		segments := make([]definition.Segments, 0)
		for _, atom := range atoms[segment.Start:segment.End] {
			segments = append(segments, atom.TextSelector)
		}
		text = atoms[0].Text
		textSelector = definition.JoinSegments(segments...)
	}

	return ParsingNode{
		Atom: definition.Atom{
			Symbol:       symbol,
			Attributes:   attributes,
			Text:         text,
			TextSelector: textSelector,
		},
		Segment:  segment,
		Children: nil,
	}
}

func stringParsingNode(node *ParsingNode, indent int, builder *strings.Builder) {
	builder.WriteString(strings.Repeat(" ", indent))
	builder.WriteString(fmt.Sprintf("%v[%v..%v): '%v'", node.Atom.Symbol, node.Segment.Start, node.Segment.End, string(node.Atom.SelectText())))
	builder.WriteString("\n")
	for _, child := range node.Children {
		stringParsingNode(child, indent+2, builder)
	}
}

func StringParsingNode(node *ParsingNode) string {
	var builder strings.Builder
	stringParsingNode(node, 0, &builder)
	return builder.String()
}
