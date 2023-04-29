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

const (
	white = 0
	gray  = 1
	black = 2
)

func orderFrom(v string, graph map[string][]string, color map[string]int, parent map[string]string, order *[]string) error {
	color[v] = gray
	for _, u := range graph[v] {
		if color[u] == white {
			parent[u] = v
			if err := orderFrom(u, graph, color, parent, order); err != nil {
				return err
			}
		} else if color[u] == gray {
			cycle := []string{v}
			for cycle[len(cycle)-1] != u {
				cycle = append(cycle, parent[cycle[len(cycle)-1]])
			}
			return fmt.Errorf("possible rules cycle detected: %v", cycle)
		}
	}
	color[v] = black
	*order = append(*order, v)
	return nil
}

func order(graph map[string][]string) ([]string, map[string]int, error) {
	color := make(map[string]int)
	parent := make(map[string]string)
	order := make([]string, 0, len(graph))
	for v := range graph {
		if color[v] == black {
			continue
		}
		err := orderFrom(v, graph, color, parent, &order)
		if err != nil {
			return nil, nil, err
		}
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
				return nil, nil, fmt.Errorf("cycle detected: incorrect order between rules %v and %v", v, next)
			}
		}
	}
	return order, position, nil
}

func (rs Rules) order() ([]string, map[string]int, error) {
	if !rs.normalized() {
		panic(fmt.Errorf("rules must be normalized"))
	}
	empty := rs.empty()
	graph := make(map[string][]string)
	for _, r := range rs {
		graph[r.Name] = make([]string, 0)
	}
	for _, r := range rs {
		switch peg := r.Expr.(type) {
		case terminals:
			continue
		case symbol:
			add(graph, r.Name, peg.Symbol)
		case junction:
			for _, j := range peg.Exprs {
				if s, ok := j.(symbol); ok {
					add(graph, r.Name, s.Symbol)
					if !empty[s.Symbol] {
						break
					}
				} else {
					break
				}
			}
		case choice:
			for _, c := range peg.Exprs {
				if s, ok := c.(symbol); ok {
					add(graph, r.Name, s.Symbol)
				}
			}
		case negation:
			if s, ok := peg.Expr.(symbol); ok {
				add(graph, r.Name, s.Symbol)
			}
		case kleene:
			if s, ok := peg.Expr.(symbol); ok {
				add(graph, r.Name, s.Symbol)
			}
		default:
			panic(fmt.Errorf("unexpected peg expression type: %v", r.Expr))
		}
	}
	return order(graph)
}
