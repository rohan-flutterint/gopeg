package parser

func isEmpty(e Expr, rs map[string]Rule, n map[string]bool) bool {
	switch peg := e.(type) {
	case junction:
		empty := true
		for _, j := range peg.Exprs {
			empty = empty && isEmpty(j, rs, n)
		}
		return empty
	case choice:
		empty := false
		for _, c := range peg.Exprs {
			empty = empty || isEmpty(c, rs, n)
		}
		return empty
	case symbol:
		return checkRule(rs[peg.Symbol], rs, n)
	case dot:
		return false
	case byteSequence:
		return len(peg.value) == 0
	case attrSequenceMatcher:
		return len(peg) == 0
	default:
		return true
	}
}

func checkRule(r Rule, rs map[string]Rule, empty map[string]bool) bool {
	value, ok := empty[r.Name]
	if ok {
		return value
	}
	empty[r.Name] = true
	value = isEmpty(r.Expr, rs, empty)
	empty[r.Name] = value
	return value
}

func (rs Rules) empty() map[string]bool {
	empty := make(map[string]bool)
	m := make(map[string]Rule)
	for _, r := range rs {
		m[r.Name] = r
	}
	for _, r := range rs {
		checkRule(r, m, empty)
	}
	//fmt.Printf("empty: %v\n", empty)
	return empty
}
