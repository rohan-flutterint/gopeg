package parser

import (
	"fmt"
)

func add(graph map[string][]string, a, b string) {
	if _, ok := graph[a]; !ok {
		graph[a] = make([]string, 0)
	}
	graph[a] = append(graph[a], b)
}

func orderFrom(v string, graph map[string][]string, used map[string]bool, order *[]string) {
	used[v] = true
	for _, u := range graph[v] {
		if used[u] {
			continue
		}
		orderFrom(u, graph, used, order)
	}
	*order = append(*order, v)
}

func order(graph map[string][]string) ([]string, map[string]int, error) {
	used := make(map[string]bool)
	order := make([]string, 0, len(graph))
	for v := range graph {
		if used[v] {
			continue
		}
		orderFrom(v, graph, used, &order)
	}
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	position := make(map[string]int)
	for i, v := range order {
		position[v] = i
	}
	for _, v := range order {
		for _, next := range graph[v] {
			if position[next] <= position[v] {
				return nil, nil, fmt.Errorf("detect cycle: incorrect order between rules %v and %v", v, next)
			}
		}
	}
	return order, position, nil
}

func (rs Rules) order() ([]string, map[string]int, error) {
	if !rs.normalized() {
		panic(fmt.Errorf("rules must be normalized"))
	}
	nonempty := rs.nonempty()
	graph := make(map[string][]string)
	for _, r := range rs {
		switch peg := r.Expr.(type) {
		case terminals:
			continue
		case nonterminal:
			add(graph, r.Name, peg.NonTerminal)
		case junction:
			for _, j := range peg.Exprs {
				if nt, ok := j.(nonterminal); ok {
					add(graph, r.Name, nt.NonTerminal)
					if nonempty[nt.NonTerminal] {
						break
					}
				} else {
					break
				}
			}
		case choice:
			for _, c := range peg.Exprs {
				if nt, ok := c.(nonterminal); ok {
					add(graph, r.Name, nt.NonTerminal)
				}
			}
		case negation:
			if nt, ok := peg.Expr.(nonterminal); ok {
				add(graph, r.Name, nt.NonTerminal)
			}
		case repetition:
			if nt, ok := peg.Expr.(nonterminal); ok {
				add(graph, r.Name, nt.NonTerminal)
			}
		default:
			panic(fmt.Errorf("unexpected peg expression type: %v", r.Expr))
		}
	}
	return order(graph)
}
