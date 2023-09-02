package parser

import (
	"errors"
	"fmt"
	"gopeg/analysis"
	"gopeg/definition"
	"strings"
)

type step struct {
	ok      bool
	advance int
}

var (
	TextNotMatchErr = errors.New("TextNotMatch")
)

func ParseAtoms(rules definition.Rules, root string, atoms []definition.Atom) (*ParsingNode, error) {
	terminalsType, err := analysis.CheckRulesConsistency(rules)
	if err != nil {
		return nil, fmt.Errorf("rules must be consistent: %w", err)
	}
	if terminalsType != analysis.AnyTerminalType && terminalsType != analysis.AtomTerminalType {
		return nil, fmt.Errorf("rules must be compatible with %v", analysis.AtomTerminalType)
	}
	rules = analysis.DesugarRules(rules)
	rules, transformation := analysis.NormalizeRules(rules)
	order, position, err := OrderRules(rules)
	if err != nil {
		return nil, fmt.Errorf("unable to topologically order rules: %w", err)
	}
	ruleMap := buildRuleMap(rules)
	table := buildStepTable(ruleMap, order, position, atoms)
	derivation, err := buildDerivationTree(ruleMap, transformation.Forward[root], position, table, atoms)
	if err != nil {
		return nil, err
	}
	parsing := transform(transformation.Backward, derivation)
	if len(parsing) != 1 {
		return nil, fmt.Errorf("tree with multiple root was formed")
	}
	return parsing[0], nil
}

func ParseText(rules definition.Rules, root string, text []byte) (*ParsingNode, error) {
	terminalsType, err := analysis.CheckRulesConsistency(rules)
	if err != nil {
		return nil, fmt.Errorf("rules must be consistent: %w", err)
	}
	if terminalsType != analysis.AnyTerminalType && terminalsType != analysis.ByteTerminalType {
		return nil, fmt.Errorf("rules must be compatible with %v", analysis.ByteTerminalType)
	}
	rules = analysis.DesugarRules(rules)
	rules, transformation := analysis.NormalizeRules(rules)
	order, position, err := OrderRules(rules)
	if err != nil {
		return nil, fmt.Errorf("unable to topologically order rules: %w", err)
	}
	ruleMap := buildRuleMap(rules)
	table := buildStepTable(ruleMap, order, position, text)
	derivation, err := buildDerivationTree(ruleMap, transformation.Forward[root], position, table, text)
	if err != nil {
		return nil, err
	}
	parsing := transform(transformation.Backward, derivation)
	if len(parsing) != 1 {
		return nil, fmt.Errorf("tree with multiple root was formed")
	}
	return parsing[0], nil
}

func advance[T any](i int, expr definition.Expr, position map[string]int, table [][]step, data []T) step {
	switch peg := expr.(type) {
	case definition.Terminals:
		advance, ok := definition.Accept[T](peg, data[i:])
		return step{ok: ok, advance: advance}
	case definition.Symbol:
		return table[i][position[peg.Name]]
	default:
		panic(fmt.Errorf("invalid usage of advance: unexpected peg expression type: %#v", expr))
	}
}

func buildDerivationTree[T any](
	ruleMap map[string]definition.Rule,
	root string,
	position map[string]int,
	table [][]step,
	data []T,
) (*ParsingNode, error) {
	if !table[0][position[root]].ok {
		return nil, TextNotMatchErr
	}
	rootNode := NewParsingNode[T](
		root,
		nil,
		data,
		definition.Segment{Start: 0, End: table[0][position[root]].advance},
	)
	derivation := []*ParsingNode{&rootNode}
	for i := 0; i < len(derivation); i++ {
		current := derivation[i]
		switch peg := ruleMap[current.Atom.Symbol].Expr.(type) {
		case definition.Terminals:
			continue
		case definition.Negation:
			continue
		case definition.Symbol:
			next := NewParsingNode[T](peg.Name, nil, data, current.Segment)
			current.Children = append(current.Children, &next)
			derivation = append(derivation, &next)
			continue
		case definition.Kleene:
			p := current.Segment.Start
			for {
				step := advance(p, peg.Expr, position, table, data)
				if !step.ok || step.advance == 0 {
					break
				}
				if s, ok := peg.Expr.(definition.Symbol); ok {
					next := NewParsingNode[T](s.Name, nil, data, definition.Segment{Start: p, End: p + step.advance})
					current.Children = append(current.Children, &next)
					derivation = append(derivation, &next)
				}
				p += step.advance
			}
		case definition.Junction:
			p := current.Segment.Start
			for _, j := range peg.Exprs {
				step := advance(p, j, position, table, data)
				if s, ok := j.(definition.Symbol); ok {
					next := NewParsingNode[T](s.Name, nil, data, definition.Segment{Start: p, End: p + step.advance})
					current.Children = append(current.Children, &next)
					derivation = append(derivation, &next)
				}
				p += step.advance
			}
		case definition.Choice:
			for _, c := range peg.Exprs {
				step := advance(current.Segment.Start, c, position, table, data)
				if !step.ok {
					continue
				}
				if s, ok := c.(definition.Symbol); ok {
					next := NewParsingNode[T](s.Name, nil, data, definition.Segment{Start: current.Segment.Start, End: current.Segment.Start + step.advance})
					current.Children = append(current.Children, &next)
					derivation = append(derivation, &next)
				}
				break
			}
		}
	}
	return &rootNode, nil
}

func buildStepTable[T any](ruleMap map[string]definition.Rule, order []string, position map[string]int, data []T) [][]step {
	// todo (sivukhin, 2023-09-02): should we use single table of size (len(text)+1) * len(order) in order to reduce amount of allocations and GC pressure?
	table := make([][]step, len(data)+1)
	for i := 0; i <= len(data); i++ {
		table[i] = make([]step, len(order))
	}

	for i := len(data); i >= 0; i-- {
		for s := len(order) - 1; s >= 0; s-- {
			switch peg := ruleMap[order[s]].Expr.(type) {
			case definition.Terminals:
				table[i][s] = advance(i, peg, position, table, data)
			case definition.Symbol:
				table[i][s] = advance(i, peg, position, table, data)
			case definition.Kleene:
				next := advance(i, peg.Expr, position, table, data)
				if next.ok && next.advance > 0 {
					table[i][s] = step{ok: true, advance: table[i+next.advance][s].advance + next.advance}
				} else {
					table[i][s] = step{ok: true, advance: 0}
				}
			case definition.Junction:
				current := i
				ok := true
				for _, j := range peg.Exprs {
					next := advance(current, j, position, table, data)
					if !next.ok {
						ok = false
					} else {
						current += next.advance
					}
				}
				if ok {
					table[i][s] = step{ok: true, advance: current - i}
				}
			case definition.Choice:
				for _, c := range peg.Exprs {
					next := advance(i, c, position, table, data)
					if next.ok {
						table[i][s] = next
						break
					}
				}
			case definition.Negation:
				next := advance(i, peg.Expr, position, table, data)
				if !next.ok {
					table[i][s] = step{ok: true, advance: 0}
				}
			}
		}
	}
	return table
}

func transform(mapping map[string]string, node *ParsingNode) []*ParsingNode {
	if name, ok := mapping[node.Atom.Symbol]; !strings.HasPrefix(name, "#") && ok {
		atom := node.Atom
		atom.Symbol = name
		next := ParsingNode{
			Atom:     atom,
			Segment:  node.Segment,
			Children: nil,
		}
		for _, child := range node.Children {
			next.Children = append(next.Children, transform(mapping, child)...)
		}
		return []*ParsingNode{&next}
	}
	nodes := make([]*ParsingNode, 0, len(node.Children))
	for _, child := range node.Children {
		nodes = append(nodes, transform(mapping, child)...)
	}
	return nodes
}
