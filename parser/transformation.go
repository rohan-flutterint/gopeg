package parser

//
//import (
//	"gopeg/definition"
//	"strings"
//)
//
//type Mapping struct {
//	id bool
//	m  map[string]string
//}
//
//type Transformation struct{ Forward, Backward Mapping }
//
//func NewIdTransformation() Transformation {
//	return Transformation{Forward: Mapping{id: true}, Backward: Mapping{id: true}}
//}
//
//func NewTransformation(forward map[string]string) Transformation {
//	backward := make(map[string]string, len(forward))
//	for a, b := range forward {
//		backward[b] = a
//	}
//	return Transformation{Forward: Mapping{m: forward}, Backward: Mapping{m: backward}}
//}
//
//func (t Mapping) get(s string) string {
//	if t.id {
//		return s
//	}
//	return t.m[s]
//}
//
//func (t Mapping) transform(n definition.ParsingNode) []definition.ParsingNode {
//	if t.id {
//		return []definition.ParsingNode{n}
//	}
//	symbol := n.Symbol()
//	if name, ok := t.m[symbol]; !strings.HasPrefix(symbol, "#") && ok {
//		start, end := n.Range()
//		attrs := make(map[string]any)
//		for key, value := range n.Attrs() {
//			if key == symbol {
//				attrs[name] = value
//			} else {
//				attrs[key] = value
//			}
//		}
//		next := &definition.parsingNode{symbol: name, start: start, end: end, children: make([]definition.ParsingNode, 0), attrs: attrs}
//		for _, child := range n.Children() {
//			next.children = append(next.children, t.transform(child)...)
//		}
//		return []definition.ParsingNode{next}
//	}
//	children := n.Children()
//	nodes := make([]definition.ParsingNode, 0, len(children))
//	for _, child := range children {
//		nodes = append(nodes, t.transform(child)...)
//	}
//	return nodes
//}
