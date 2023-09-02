package parser

import (
	"fmt"
	"gopeg/definition"
)

func addBackwardDeps(deps map[string][]string, root string, exprs []definition.Expr) {
	for _, expr := range exprs {
		if symbol, ok := expr.(definition.Symbol); ok {
			deps[symbol.Name] = append(deps[symbol.Name], root)
		}
	}
}

func selectForwardDeps(exprs []definition.Expr) []string {
	deps := make([]string, 0)
	for _, expr := range exprs {
		if symbol, ok := expr.(definition.Symbol); ok {
			deps = append(deps, symbol.Name)
		}
	}
	return deps
}

func buildRuleMap(rules definition.Rules) map[string]definition.Rule {
	ruleMap := make(map[string]definition.Rule)
	for _, rule := range rules {
		ruleMap[rule.Name] = rule
	}
	return ruleMap
}

func buildForwardRuleDeps(rules definition.Rules) map[string][]string {
	ruleDeps := make(map[string][]string)
	for _, rule := range rules {
		switch peg := rule.Expr.(type) {
		case definition.Junction:
			ruleDeps[rule.Name] = selectForwardDeps(peg.Exprs)
		case definition.Choice:
			ruleDeps[rule.Name] = selectForwardDeps(peg.Exprs)
		case definition.Kleene:
			ruleDeps[rule.Name] = selectForwardDeps([]definition.Expr{peg.Expr})
		case definition.Negation:
			ruleDeps[rule.Name] = selectForwardDeps([]definition.Expr{peg.Expr})
		case definition.Symbol:
			ruleDeps[rule.Name] = selectForwardDeps([]definition.Expr{rule.Expr})
		case definition.Terminals:
			continue
		default:
			panic(fmt.Errorf("unexpected peg expression type: %#v", rule.Expr))
		}
	}
	return ruleDeps
}

func buildBackwardRuleDeps(rules definition.Rules) map[string][]string {
	ruleDeps := make(map[string][]string)
	for _, rule := range rules {
		ruleDeps[rule.Name] = make([]string, 0)
	}
	for _, rule := range rules {
		switch peg := rule.Expr.(type) {
		case definition.Junction:
			addBackwardDeps(ruleDeps, rule.Name, peg.Exprs)
		case definition.Choice:
			addBackwardDeps(ruleDeps, rule.Name, peg.Exprs)
		case definition.Kleene:
			addBackwardDeps(ruleDeps, rule.Name, []definition.Expr{peg.Expr})
		case definition.Negation:
			addBackwardDeps(ruleDeps, rule.Name, []definition.Expr{peg.Expr})
		case definition.Symbol:
			addBackwardDeps(ruleDeps, rule.Name, []definition.Expr{rule.Expr})
		case definition.Terminals:
			continue
		default:
			panic(fmt.Errorf("unexpected peg expression type: %#v", rule.Expr))
		}
	}
	return ruleDeps
}
