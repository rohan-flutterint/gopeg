package parser

import (
	"fmt"
)

func desugarMany(exprs []Expr) []Expr {
	result := make([]Expr, 0, len(exprs))
	for _, c := range exprs {
		result = append(result, desugar(c))
	}
	return result
}

func desugar(expr Expr) Expr {
	switch peg := expr.(type) {
	case optional:
		return choice{Exprs: []Expr{desugar(peg.Expr), NewEmpty()}}
	case ensure:
		return negation{negation{Expr: desugar(peg.Expr)}}
	case repetition:
		exprs := make([]Expr, 0, peg.min+1)
		for i := 0; i < peg.min; i++ {
			exprs = append(exprs, peg.Expr)
		}
		exprs = append(exprs, kleene{peg.Expr})
		return junction{exprs}
	case terminals:
		return peg
	case symbol:
		return peg
	case junction:
		return junction{desugarMany(peg.Exprs)}
	case choice:
		return choice{desugarMany(peg.Exprs)}
	case kleene:
		return kleene{desugar(peg.Expr)}
	case negation:
		return negation{desugar(peg.Expr)}
	default:
		panic(fmt.Errorf("unexpected peg expression type: %v", expr))
	}
}

func (r Rule) desugar() Rule {
	return Rule{Name: r.Name, Expr: desugar(r.Expr)}
}

func (rs Rules) desugar() Rules {
	result := make(Rules, 0, len(rs))
	for _, r := range rs {
		result = append(result, r.desugar())
	}
	return result
}
