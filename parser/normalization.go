package parser

import "fmt"

func createName(name string, id int) string { return fmt.Sprintf("%v#%v", name, id) }

func leaf(e Expr) (Expr, bool) {
	switch peg := e.(type) {
	case terminals:
		return peg, true
	case symbol:
		return symbol{Symbol: createName(peg.Symbol, 0)}, true
	}
	return nil, false
}

func normalize(expr Expr, name string, id int) Rules {
	rootName := createName(name, id)
	switch peg := expr.(type) {
	case terminals:
		return Rules{{rootName, peg}}
	case symbol:
		return Rules{{rootName, symbol{Symbol: createName(peg.Symbol, 0)}}}
	case kleene:
		if l, ok := leaf(peg.Expr); ok {
			return Rules{Rule{rootName, kleene{l}}}
		}
		root := Rule{rootName, kleene{symbol{createName(name, id+1)}}}
		return append(Rules{root}, normalize(peg.Expr, name, id+1)...)
	case negation:
		if l, ok := leaf(peg.Expr); ok {
			return Rules{Rule{rootName, negation{l}}}
		}
		root := Rule{rootName, negation{symbol{createName(name, id+1)}}}
		return append(Rules{root}, normalize(peg.Expr, name, id+1)...)
	case junction:
		js := Rules{}
		var roots []Expr
		for _, j := range peg.Exprs {
			if l, ok := leaf(j); ok {
				roots = append(roots, l)
			} else {
				roots = append(roots, symbol{createName(name, id+1+len(js))})
				js = append(js, normalize(j, name, id+1+len(js))...)
			}
		}
		root := Rule{rootName, junction{roots}}
		return append(Rules{root}, js...)
	case choice:
		cs := Rules{}
		var roots []Expr
		for _, c := range peg.Exprs {
			if l, ok := leaf(c); ok {
				roots = append(roots, l)
			} else {
				roots = append(roots, symbol{createName(name, id+1+len(cs))})
				cs = append(cs, normalize(c, name, id+1+len(cs))...)
			}
		}
		root := Rule{rootName, choice{roots}}
		return append(Rules{root}, cs...)
	default:
		panic(fmt.Errorf("unexpected peg expression type: %v", expr))
	}
}

func (rs Rules) normalize() (Rules, Transformation) {
	rules := make(Rules, 0)
	mapping := make(map[string]string)
	for _, r := range rs {
		mapping[r.Name] = createName(r.Name, 0)
		rules = append(rules, normalize(r.Expr, r.Name, 0)...)
	}
	return rules, NewTransformation(mapping)
}

func leafs(es []Expr) bool {
	for _, e := range es {
		if _, ok := leaf(e); !ok {
			return false
		}
	}
	return true
}

func (r Rule) normalized() bool {
	switch peg := r.Expr.(type) {
	case terminals:
		return true
	case symbol:
		return true
	case kleene:
		_, ok := leaf(peg.Expr)
		return ok
	case negation:
		_, ok := leaf(peg.Expr)
		return ok
	case junction:
		return leafs(peg.Exprs)
	case choice:
		return leafs(peg.Exprs)
	default:
		panic(fmt.Errorf("unexpected peg expression type: %v", r.Expr))
	}
}

func (rs Rules) normalized() bool {
	for _, r := range rs {
		if !r.normalized() {
			return false
		}
	}
	return true
}
