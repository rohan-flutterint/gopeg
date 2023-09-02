package analysis

import (
	"fmt"
	"gopeg/definition"
)

func CheckDesugaredExprs(expr []definition.Expr) error {
	for _, e := range expr {
		if err := CheckDesugaredExpr(e); err != nil {
			return err
		}
	}
	return nil
}

func CheckDesugaredExpr(expr definition.Expr) error {
	switch peg := expr.(type) {
	case definition.Optional:
		return fmt.Errorf("found Optional expression")
	case definition.Ensure:
		return fmt.Errorf("found Ensure expression")
	case definition.Repetition:
		return fmt.Errorf("found Repetition expression")
	case definition.Junction:
		return CheckDesugaredExprs(peg.Exprs)
	case definition.Choice:
		return CheckDesugaredExprs(peg.Exprs)
	case definition.Kleene:
		return CheckDesugaredExpr(peg.Expr)
	case definition.Negation:
		return CheckDesugaredExpr(peg.Expr)
	default:
		return nil
	}
}

func CheckDesugaredRule(rule definition.Rule) error {
	return CheckDesugaredExpr(rule.Expr)
}

func CheckDesugaredRules(rules definition.Rules) error {
	for _, rule := range rules {
		if err := CheckDesugaredRule(rule); err != nil {
			return fmt.Errorf("non-desugared expression in rule '%v': %w", rule.Name, err)
		}
	}
	return nil
}

func DesugarExprs(expr []definition.Expr) []definition.Expr {
	desugared := make([]definition.Expr, 0, len(expr))
	for _, e := range expr {
		desugared = append(desugared, DesugarExpr(e))
	}
	return desugared
}

func DesugarExpr(expr definition.Expr) definition.ExprCore {
	switch peg := expr.(type) {
	case definition.Optional:
		return definition.Choice{Exprs: []definition.Expr{DesugarExpr(peg.Expr), definition.NewEmpty()}}
	case definition.Ensure:
		return definition.Negation{Expr: definition.Negation{Expr: DesugarExpr(peg.Expr)}}
	case definition.Repetition:
		desugared := DesugarExpr(peg.Expr)
		exprs := make([]definition.Expr, 0, int(peg.Min)+1)
		for i := 0; i < int(peg.Min); i++ {
			exprs = append(exprs, desugared)
		}
		return definition.Junction{Exprs: append(exprs, definition.Kleene{Expr: desugared})}
	case definition.Junction:
		return definition.Junction{Exprs: DesugarExprs(peg.Exprs)}
	case definition.Choice:
		return definition.Choice{Exprs: DesugarExprs(peg.Exprs)}
	case definition.Kleene:
		return definition.Kleene{Expr: DesugarExpr(peg.Expr)}
	case definition.Negation:
		return definition.Negation{Expr: DesugarExpr(peg.Expr)}
	case definition.ExprCore:
		return peg
	default:
		panic(fmt.Errorf("unexpected peg expression type: %v", expr))
	}
}

func DesugarRule(rule definition.Rule) definition.Rule {
	return definition.Rule{Name: rule.Name, Expr: DesugarExpr(rule.Expr)}
}

func DesugarRules(rules definition.Rules) definition.Rules {
	result := make(definition.Rules, 0, len(rules))
	for _, rule := range rules {
		result = append(result, DesugarRule(rule))
	}
	return result
}
