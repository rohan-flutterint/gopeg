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

func StringParsingNode(n ParsingNode, text string) string {
	var printNode func(*strings.Builder, ParsingNode, int)
	printNode = func(b *strings.Builder, pn ParsingNode, offset int) {
		start, end := pn.Range()
		if pn.NonTerminal() == "" {
			b.WriteString(fmt.Sprintf("%v'%v'\n", strings.Repeat(" ", offset), text[start:end]))
		} else {
			b.WriteString(fmt.Sprintf("%v%v: [%v..%v)\n", strings.Repeat(" ", offset), pn.NonTerminal(), start, end))
		}
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

func (rs Rules) parse(root, text string) (ParsingNode, error) {
	n, mapping := rs, NewIdTransformation()
	if !rs.normalized() {
		n, mapping = rs.normalize()
	}
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
		case nonterminal:
			return table[i][position[peg.NonTerminal]]
		default:
			panic(fmt.Errorf("invalid usage of advance: unexpected peg expression type: %#v", expr))
		}
	}
	for i := len(text); i >= 0; i-- {
		for s := len(order) - 1; s >= 0; s-- {
			switch peg := access[order[s]].Expr.(type) {
			case terminals:
				table[i][s] = advance(i, peg)
			case nonterminal:
				table[i][s] = advance(i, peg)
			case repetition:
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
		}
	}
	rootN := mapping.Forward.get(root)
	if !table[0][position[rootN]].ok {
		return nil, TextNotMatchErr
	}
	derivation := []*parsingNode{{nonTerminal: rootN, start: 0, end: table[0][position[rootN]].advance}}
	for i := 0; i < len(derivation); i++ {
		current := derivation[i]
		if current.nonTerminal == "" {
			continue
		}
		switch peg := access[current.nonTerminal].Expr.(type) {
		case terminals:
			next := &parsingNode{nonTerminal: "", start: current.start, end: current.end}
			current.children = append(current.children, next)
			continue
		case negation:
			continue
		case nonterminal:
			next := &parsingNode{nonTerminal: peg.NonTerminal, start: current.start, end: current.end}
			current.children = append(current.children, next)
			derivation = append(derivation, next)
			continue
		case repetition:
			p := current.start
			for {
				step := advance(p, peg.Expr)
				if !step.ok {
					break
				}
				nt, _ := NonTerminalName(peg.Expr)
				next := &parsingNode{nonTerminal: nt, start: p, end: p + step.advance}
				current.children = append(current.children, next)
				derivation = append(derivation, next)
				p += step.advance
			}
		case junction:
			p := current.start
			for _, j := range peg.Exprs {
				step := advance(p, j)
				if step.advance > 0 {
					nt, _ := NonTerminalName(j)
					next := &parsingNode{nonTerminal: nt, start: p, end: p + step.advance}
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
				if step.advance > 0 {
					nt, _ := NonTerminalName(c)
					next := &parsingNode{nonTerminal: nt, start: current.start, end: current.start + step.advance}
					current.children = append(current.children, next)
					derivation = append(derivation, next)
				}
				break
			}
		}
	}
	return mapping.Backward.transform(derivation[0])[0], nil
}
