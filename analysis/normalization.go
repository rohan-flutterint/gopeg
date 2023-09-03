package analysis

import (
	"fmt"
	"gopeg/definition"
)

type nameGenerator struct {
	name string
	id   int
}

func (g *nameGenerator) createRootName(name string) string {
	return fmt.Sprintf("%v#%v", name, 0)
}
func (g *nameGenerator) createNextName() string {
	uniqueName := fmt.Sprintf("%v#%v", g.name, g.id)
	g.id++
	return uniqueName
}

func prepareExpr(generator *nameGenerator, expr definition.Expr) (definition.Expr, definition.Rules) {
	normalized, rules := normalizeExpr(generator, expr)
	if isLeaf(normalized) {
		return normalized, rules
	}
	symbol := definition.Symbol{Name: generator.createNextName()}
	return symbol, append(rules, definition.Rule{Name: symbol.Name, Expr: normalized})
}

func normalizeExpr(generator *nameGenerator, expr definition.Expr) (definition.Expr, definition.Rules) {
	switch peg := expr.(type) {
	case definition.Terminals:
		return peg, nil
	case definition.Symbol:
		return definition.Symbol{Name: generator.createRootName(peg.Name), Attributes: peg.Attributes}, nil
	case definition.Kleene:
		normalized, rules := prepareExpr(generator, peg.Expr)
		return definition.Kleene{Expr: normalized}, rules
	case definition.Negation:
		normalized, rules := prepareExpr(generator, peg.Expr)
		return definition.Negation{Expr: normalized}, rules
	case definition.Junction:
		var rules definition.Rules
		children := make([]definition.Expr, 0, len(peg.Exprs))
		for _, e := range peg.Exprs {
			normalized, rulesE := prepareExpr(generator, e)
			rules = append(rules, rulesE...)
			children = append(children, normalized)
		}
		return definition.Junction{Exprs: children}, rules
	case definition.Choice:
		var rules definition.Rules
		children := make([]definition.Expr, 0, len(peg.Exprs))
		for _, e := range peg.Exprs {
			normalized, rulesE := prepareExpr(generator, e)
			rules = append(rules, rulesE...)
			children = append(children, normalized)
		}
		return definition.Choice{Exprs: children}, rules
	default:
		panic(fmt.Errorf("unexpected peg expression type: %#v", expr))
	}
}

func NormalizeRules(rules definition.Rules) (definition.Rules, Transformation) {
	normalized := make(definition.Rules, 0, len(rules))
	mapping := make(map[string]string)
	rulesByName := make(map[string][]definition.Expr)
	ruleNames := make([]string, 0)
	for _, rule := range rules {
		if _, ok := rulesByName[rule.Name]; !ok {
			rulesByName[rule.Name] = make([]definition.Expr, 0)
			ruleNames = append(ruleNames, rule.Name)
		}
		rulesByName[rule.Name] = append(rulesByName[rule.Name], rule.Expr)
	}
	for _, ruleName := range ruleNames {
		generator := nameGenerator{name: ruleName, id: 0}
		rootName := generator.createNextName()
		mapping[ruleName] = rootName
		normalizedExpr, additionalRules := normalizeExpr(&generator, definition.NewChoice(rulesByName[ruleName]...))
		normalized = append(normalized, definition.Rule{Name: rootName, Expr: normalizedExpr})
		normalized = append(normalized, additionalRules...)
	}
	return normalized, NewTransformation(mapping)
}

func CheckNormalizedRules(rules definition.Rules) error {
	ruleNames := make(map[string]struct{})
	for _, rule := range rules {
		if err := CheckNormalizedRule(rule); err != nil {
			return fmt.Errorf("non-normalized expression in rule '%v': %w", rule.Name, err)
		}
		if _, ok := ruleNames[rule.Name]; ok {
			return fmt.Errorf("duplicate rule names: '%v'", rule.Name)
		}
		ruleNames[rule.Name] = struct{}{}
	}
	return nil
}

func CheckNormalizedRule(rule definition.Rule) error {
	switch peg := rule.Expr.(type) {
	case definition.Terminals:
		return nil
	case definition.Symbol:
		return nil
	case definition.Kleene:
		return checkLeaf(peg.Expr)
	case definition.Negation:
		return checkLeaf(peg.Expr)
	case definition.Junction:
		return checkLeafs(peg.Exprs)
	case definition.Choice:
		return checkLeafs(peg.Exprs)
	default:
		panic(fmt.Errorf("unexpected peg expression type: %#v", rule.Expr))
	}
}

func isType[T any](v any) bool {
	_, ok := v.(T)
	return ok
}

func isLeaf(expr definition.Expr) bool {
	return isType[definition.Terminals](expr) || isType[definition.Symbol](expr)
}

func checkLeaf(expr definition.Expr) error {
	if !isLeaf(expr) {
		return fmt.Errorf("expression must be leaf: %v", expr)
	}
	return nil
}

func checkLeafs(expr []definition.Expr) error {
	for _, e := range expr {
		if err := checkLeaf(e); err != nil {
			return err
		}
	}
	return nil
}
