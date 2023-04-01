package parser

type Mapping struct {
	id bool
	m  map[string]string
}

type Transformation struct{ Forward, Backward Mapping }

func NewIdTransformation() Transformation {
	return Transformation{Forward: Mapping{id: true}, Backward: Mapping{id: true}}
}

func NewTransformation(forward map[string]string) Transformation {
	backward := make(map[string]string, len(forward))
	for a, b := range forward {
		backward[b] = a
	}
	return Transformation{Forward: Mapping{m: forward}, Backward: Mapping{m: backward}}
}

func (t Mapping) get(s string) string {
	if t.id {
		return s
	}
	return t.m[s]
}

func (t Mapping) transform(n ParsingNode) []ParsingNode {
	if t.id || n.NonTerminal() == "" {
		return []ParsingNode{n}
	}
	if name, ok := t.m[n.NonTerminal()]; ok {
		start, end := n.Range()
		next := &parsingNode{nonTerminal: name, start: start, end: end, children: make([]ParsingNode, 0)}
		for _, child := range n.Children() {
			next.children = append(next.children, t.transform(child)...)
		}
		return []ParsingNode{next}
	}
	children := n.Children()
	nodes := make([]ParsingNode, 0, len(children))
	for _, child := range children {
		nodes = append(nodes, t.transform(child)...)
	}
	return nodes
}
