package parser

import (
	"errors"
	"fmt"
	"strings"
)

type step struct {
	ok      bool
	advance int
}

func StringParsingNode[T any](n ParsingNode, text []T) string {
	var printNode func(*strings.Builder, ParsingNode, int)
	printNode = func(b *strings.Builder, pn ParsingNode, offset int) {
		start, end := pn.Range()
		b.WriteString(fmt.Sprintf("%v%v: [%v..%v)\n", strings.Repeat(" ", offset), pn.Symbol(), start, end))
		for _, c := range pn.Children() {
			printNode(b, c, offset+2)
		}
		return
	}
	b := strings.Builder{}
	printNode(&b, n, 0)
	return b.String()
}

var (
	TextNotMatchErr = errors.New("TextNotMatch")
)

func Parse[T any](rs Rules, root string, text []T) (ParsingNode, error) {
	rs = rs.desugar()
	n, mapping := rs, NewIdTransformation()
	if !rs.normalized() {
		n, mapping = rs.normalize()
	}
	//fmt.Printf("n: %v\n", n)
	access := make(map[string]Rule)
	for _, r := range n {
		access[r.Name] = r
	}

	order, position, err := n.order()
	if err != nil {
		return nil, err
	}

	table := make([][]step, len(text)+1)
	for i := 0; i <= len(text); i++ {
		table[i] = make([]step, len(order))
	}

	advance := func(i int, expr Expr) step {
		switch peg := expr.(type) {
		case terminals:
			advance, ok := peg.Accept(text[i:])
			return step{ok: ok, advance: advance}
		case symbol:
			return table[i][position[peg.Symbol]]
		default:
			panic(fmt.Errorf("invalid usage of advance: unexpected peg expression type: %#v", expr))
		}
	}

	for i := len(text); i >= 0; i-- {
		for s := len(order) - 1; s >= 0; s-- {
			switch peg := access[order[s]].Expr.(type) {
			case terminals:
				table[i][s] = advance(i, peg)
			case symbol:
				table[i][s] = advance(i, peg)
			case kleene:
				next := advance(i, peg.Expr)
				if next.ok && next.advance > 0 {
					table[i][s] = step{ok: true, advance: table[i+next.advance][s].advance + next.advance}
				} else {
					table[i][s] = step{ok: true, advance: 0}
				}
			case junction:
				current := i
				ok := true
				for _, j := range peg.Exprs {
					next := advance(current, j)
					if !next.ok {
						ok = false
					} else {
						current += next.advance
					}
				}
				if ok {
					table[i][s] = step{ok: true, advance: current - i}
				}
			case choice:
				for _, c := range peg.Exprs {
					next := advance(i, c)
					if next.ok {
						table[i][s] = next
						break
					}
				}
			case negation:
				next := advance(i, peg.Expr)
				if !next.ok {
					table[i][s] = step{ok: true, advance: 0}
				}
			}
			//fmt.Printf("i: %v, s: %v, %v (%v)\n", i, s, table[i][s], access[order[s]])
		}
	}

	rootN := mapping.Forward.get(root)
	//fmt.Printf("root: %v, rootN: %v %v\n", root, rootN, position[rootN])
	if !table[0][position[rootN]].ok {
		return nil, TextNotMatchErr
	}
	derivation := []*parsingNode{NewParsingNode[T](rootN, 0, table[0][position[rootN]].advance, text)}
	for i := 0; i < len(derivation); i++ {
		current := derivation[i]
		if current.symbol == "" {
			continue
		}
		switch peg := access[current.symbol].Expr.(type) {
		case terminals:
			continue
		case negation:
			continue
		case symbol:
			next := NewParsingNode[T](peg.Symbol, current.start, current.end, text)
			current.children = append(current.children, next)
			derivation = append(derivation, next)
			continue
		case kleene:
			p := current.start
			for {
				step := advance(p, peg.Expr)
				if !step.ok || step.advance == 0 {
					break
				}
				if s, ok := SymbolName(peg.Expr); ok {
					next := NewParsingNode[T](s, p, p+step.advance, text)
					current.children = append(current.children, next)
					derivation = append(derivation, next)
				}
				p += step.advance
			}
		case junction:
			p := current.start
			for _, j := range peg.Exprs {
				step := advance(p, j)
				if s, ok := SymbolName(j); ok {
					next := NewParsingNode[T](s, p, p+step.advance, text)
					current.children = append(current.children, next)
					derivation = append(derivation, next)
				}
				p += step.advance
			}
		case choice:
			for _, c := range peg.Exprs {
				step := advance(current.start, c)
				if !step.ok {
					continue
				}
				if s, ok := SymbolName(c); ok {
					next := NewParsingNode[T](s, current.start, current.start+step.advance, text)
					current.children = append(current.children, next)
					derivation = append(derivation, next)
				}
				break
			}
		}
	}
	return mapping.Backward.transform(derivation[0])[0], nil
}
