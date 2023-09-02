package parser

import (
	"fmt"
	"gopeg/analysis"
	"gopeg/definition"
)

func add(graph map[string][]string, a, b string) {
	if _, ok := graph[a]; !ok {
		graph[a] = make([]string, 0)
	}
	graph[a] = append(graph[a], b)
}

type visitState int

const (
	notVisited     visitState = 0
	visitedEntered visitState = 1
	visitedExited  visitState = 2
)

func orderFrom(v string, graph map[string][]string, color map[string]visitState, parent map[string]string, order *[]string) error {
	color[v] = visitedEntered
	for _, u := range graph[v] {
		if color[u] == notVisited {
			parent[u] = v
			if err := orderFrom(u, graph, color, parent, order); err != nil {
				return err
			}
		} else if color[u] == visitedEntered {
			cycle := []string{v}
			for cycle[len(cycle)-1] != u {
				cycle = append(cycle, parent[cycle[len(cycle)-1]])
			}
			return fmt.Errorf("possible rules cycle detected: %v", cycle)
		}
	}
	color[v] = visitedExited
	*order = append(*order, v)
	return nil
}

func topoSort(graph map[string][]string) ([]string, map[string]int, error) {
	color := make(map[string]visitState)
	parent := make(map[string]string)
	sequence := make([]string, 0, len(graph))
	for v := range graph {
		if color[v] == visitedExited {
			continue
		}
		err := orderFrom(v, graph, color, parent, &sequence)
		if err != nil {
			return nil, nil, err
		}
	}
	for i, j := 0, len(sequence)-1; i < j; i, j = i+1, j-1 {
		sequence[i], sequence[j] = sequence[j], sequence[i]
	}
	position := make(map[string]int)
	for i, v := range sequence {
		position[v] = i
	}
	for _, v := range sequence {
		for _, next := range graph[v] {
			if position[next] <= position[v] {
				return nil, nil, fmt.Errorf("cycle detected: incorrect order between rules %v and %v", v, next)
			}
		}
	}
	return sequence, position, nil
}

func OrderRules(rules definition.Rules) ([]string, map[string]int, error) {
	if _, err := analysis.CheckRulesConsistency(rules); err != nil {
		panic(fmt.Errorf("rules must be consistent: %w", err))
	}
	if err := analysis.CheckDesugaredRules(rules); err != nil {
		panic(fmt.Errorf("rules must be desugared before ordering: %w", err))
	}
	if err := analysis.CheckNormalizedRules(rules); err != nil {
		panic(fmt.Errorf("rules must be normalized before ordering: %w", err))
	}
	emptiness := GetRulesEmptiness(rules)
	ruleForwardDeps := buildForwardRuleDeps(rules)
	ruleNonEmptyDeps := make(map[string][]string)
	for _, rule := range rules {
		switch peg := rule.Expr.(type) {
		case definition.Junction:
			edges := make([]string, 0)
		depLoop:
			for _, dep := range peg.Exprs {
				switch pegDeg := dep.(type) {
				case definition.Symbol:
					edges = append(edges, pegDeg.Name)
					if !emptiness[pegDeg.Name] {
						break depLoop
					}
				case definition.Terminals:
					if !isEmptyTerminal(pegDeg) {
						break depLoop
					}
				default:
					panic(fmt.Errorf("unexpected peg expression type: %v", pegDeg))
				}
			}
			ruleNonEmptyDeps[rule.Name] = edges
		case definition.Choice:
			ruleNonEmptyDeps[rule.Name] = ruleForwardDeps[rule.Name]
		case definition.Negation:
			ruleNonEmptyDeps[rule.Name] = ruleForwardDeps[rule.Name]
		case definition.Kleene:
			ruleNonEmptyDeps[rule.Name] = ruleForwardDeps[rule.Name]
		case definition.Symbol:
			ruleNonEmptyDeps[rule.Name] = ruleForwardDeps[rule.Name]
		case definition.Terminals:
			ruleNonEmptyDeps[rule.Name] = ruleForwardDeps[rule.Name]
		default:
			panic(fmt.Errorf("unexpected peg expression type: %#v", rule.Expr))
		}
	}
	return topoSort(ruleNonEmptyDeps)
}
