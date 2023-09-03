package parser

import (
	"fmt"
	"github.com/sivukhin/gopeg/analysis"
	"github.com/sivukhin/gopeg/definition"
)

type emptinessState int

const (
	ruleIsUnknown  emptinessState = 0
	ruleIsEmpty    emptinessState = 1
	ruleIsNonEmpty emptinessState = 2
)

func isEmptyTerminal(expr definition.Terminals) bool {
	switch peg := expr.(type) {
	case definition.Dot:
		return false
	case definition.Empty:
		return true
	case definition.TextPattern:
		return peg.Regex.Match([]byte{})
	case definition.TextToken:
		return len(peg.Text) == 0
	case definition.AtomPattern:
		return false
	case definition.StartOfFile:
		return true
	case definition.EndOfFile:
		return true
	default:
		panic(fmt.Errorf("unexpected peg terminal type: %#v", expr))
	}
}

func isEmptyExpr(expr definition.Expr, emptiness map[string]emptinessState) emptinessState {
	switch peg := expr.(type) {
	case definition.Junction:
		for _, e := range peg.Exprs {
			result := isEmptyExpr(e, emptiness)
			if result == ruleIsNonEmpty {
				return ruleIsNonEmpty
			}
			if result == ruleIsUnknown {
				return ruleIsUnknown
			}
		}
		return ruleIsEmpty
	case definition.Choice:
		for _, e := range peg.Exprs {
			result := isEmptyExpr(e, emptiness)
			if result == ruleIsEmpty {
				return ruleIsEmpty
			}
			if result == ruleIsUnknown {
				return ruleIsUnknown
			}
		}
		return ruleIsNonEmpty
	case definition.Kleene:
		return ruleIsEmpty
	case definition.Negation:
		return ruleIsEmpty
	case definition.Symbol:
		return emptiness[peg.Name]
	case definition.Terminals:
		if isEmptyTerminal(peg) {
			return ruleIsEmpty
		}
		return ruleIsNonEmpty
	default:
		panic(fmt.Errorf("unexpected peg expression type: %#v", expr))
	}
}

func GetRulesEmptiness(rules definition.Rules) map[string]bool {
	if _, err := analysis.CheckRulesConsistency(rules); err != nil {
		panic(fmt.Errorf("rules must be consistent: %w", err))
	}
	if err := analysis.CheckDesugaredRules(rules); err != nil {
		panic(fmt.Errorf("rules must be desugared before checking for emptiness: %w", err))
	}
	if err := analysis.CheckNormalizedRules(rules); err != nil {
		panic(fmt.Errorf("rules must be normalized before checking for emptiness: %w", err))
	}
	ruleMap := buildRuleMap(rules)
	ruleBackwardDeps := buildBackwardRuleDeps(rules)
	ruleState := make(map[string]emptinessState)
	queue := make([]string, 0)
	for _, rule := range rules {
		result := isEmptyExpr(rule.Expr, ruleState)
		if result == ruleIsUnknown {
			continue
		}
		ruleState[rule.Name] = result
		queue = append(queue, rule.Name)
	}
	for i := 0; i < len(queue); i++ {
		current := queue[i]
		for _, dep := range ruleBackwardDeps[current] {
			if _, ok := ruleState[dep]; ok {
				continue
			}
			result := isEmptyExpr(ruleMap[dep].Expr, ruleState)
			if result == ruleIsUnknown {
				continue
			}
			ruleState[dep] = result
			queue = append(queue, dep)
		}
	}
	emptiness := make(map[string]bool)
	for _, rule := range rules {
		emptiness[rule.Name] = ruleState[rule.Name] == ruleIsEmpty
	}
	return emptiness
}
