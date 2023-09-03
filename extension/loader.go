package extension

import (
	"fmt"
	"gopeg/definition"
	"gopeg/parser"
	"strconv"
)

func Load(text string) (definition.Rules, error) {
	tokens, err := parser.ParseText(PegTokenizerRules, PegText, []byte(text))
	if err != nil {
		return nil, fmt.Errorf("unable to parse syntax structure: %w", err)
	}
	if tokens.Segment.Length() != len(text) {
		return nil, fmt.Errorf("unable to fully tokenize input (valid first %v bytes)", tokens.Segment.Length())
	}
	atoms := make([]definition.Atom, 0)
	for _, atom := range tokens.Children {
		atoms = append(atoms, atom.Atom)
	}
	peg, err := parser.ParseAtoms(PegGrammarRules, PegDefinitions, atoms)
	if err != nil {
		return nil, fmt.Errorf("unable to parse lexical structure: %w", err)
	}
	if peg.Segment.Length() != len(atoms) {
		return nil, fmt.Errorf("unable to fully parse tokenized input (valid first %v tokens)", peg.Segment.Length())
	}
	rules := make(definition.Rules, 0)
	for _, d := range peg.EnsureOnlySymbol(PegDefinition) {
		name := string(d.MustSelectBySymbol(PegName).Atom.SelectText())
		current, additional, err := rule(d.MustSelectBySymbol(PegRule), atoms)
		if err != nil {
			return nil, err
		}
		rules = append(rules, additional...)
		rules = append(rules, definition.NewRule(name, current))
	}
	return rules, nil
}

func rule(node *parser.ParsingNode, atoms []definition.Atom) (definition.Expr, []definition.Rule, error) {
	if node.Atom.Symbol != PegRule {
		panic(fmt.Errorf("unexpcted node type: %v", node.Atom.Symbol))
	}

	rules := make([]definition.Rule, 0)
	choices := make([]definition.Expr, 0, len(node.Children))
	for _, choice := range node.EnsureOnlySymbol(PegChoice) {
		junctions := make([]definition.Expr, 0, len(choice.Children))
		for _, junction := range choice.EnsureOnlySymbol(PegJunction) {
			expr := junction.MustSelectBySymbol(PegExpression)
			var current definition.Expr
			var err error
			if len(expr.Children) == 0 {
				if expr.Segment.Length() != 1 {
					panic(fmt.Errorf("unexpected expression: %#v", expr))
				}
				atom := atoms[expr.Segment.Start]
				current, err = atom2expr(atom, true)
				if err != nil {
					return nil, nil, err
				}
			} else {
				child := expr.EnsureOnlySingle()
				switch child.Atom.Symbol {
				case PegRule:
					var addition []definition.Rule
					current, addition, err = rule(child, atoms)
					if err != nil {
						return nil, nil, err
					}
					rules = append(rules, addition...)
				case PegMap:
					matcher := make(map[string]definition.TextTerminals)
					for _, keyValue := range child.EnsureOnlySymbol(PegMapKeyValue) {
						key := keyValue.MustSelectBySymbol(PegMapKey)
						value, valueOk := keyValue.TrySelectBySymbol(PegMapValue)
						if !valueOk {
							matcher[string(key.Atom.SelectText())] = nil
						} else {
							atom := atoms[value.Segment.Start]
							expr, err := atom2expr(atom, true)
							if err != nil {
								return nil, nil, err
							}
							matcher[string(key.Atom.SelectText())] = expr.(definition.TextTerminals)
						}
					}
					current = definition.NewAtomPattern(matcher)
				default:
					panic(fmt.Errorf("unexpected atom symbol: %v", child.Atom.Symbol))
				}
			}

			if prefix, ok := junction.TrySelectBySymbol(PegPrefix); ok {
				control := string(prefix.Atom.SelectText())
				switch control {
				case "!":
					current = definition.NewNegation(current)
				case "&":
					current = definition.NewEnsure(current)
				default:
					return nil, nil, fmt.Errorf("unknown prefix: %v", control)
				}
			}
			if suffix, ok := junction.TrySelectBySymbol(PegSuffix); ok {
				control := string(suffix.Atom.SelectText())
				switch control {
				case "?":
					current = definition.NewOptional(current)
				case "*":
					current = definition.NewRepetition(current)
				case "+":
					current = definition.NewRepetitionN(current, 1)
				default:
					return nil, nil, fmt.Errorf("unknown suffix: %v", control)
				}
			}
			if alias, ok := junction.TrySelectBySymbol(PegAlias); ok {
				symbol := string(alias.Atom.SelectText())
				rules = append(rules, definition.NewRule(symbol, current))
				current = definition.NewSymbol(symbol)
			}
			junctions = append(junctions, current)
		}
		choices = append(choices, definition.NewJunction(junctions...))
	}
	return definition.NewChoice(choices...), rules, nil
}

func atom2expr(atom definition.Atom, textMatcher bool) (definition.Expr, error) {
	switch atom.Symbol {
	case PegToken:
		return definition.NewSymbol(string(atom.SelectText())), nil
	case PegString:
		token := string(atom.SelectText())
		unescaped, err := strconv.Unquote(token)
		if err != nil {
			return nil, fmt.Errorf("unable to unescape %v token '%v': %w", atom.Symbol, token, err)
		}
		return definition.NewTextToken(unescaped), nil
	case PegRegex:
		token := string(atom.SelectText())
		unescaped, err := strconv.Unquote(token)
		if err != nil {
			return nil, fmt.Errorf("unable to unescape %v token '%v': %w", atom.Symbol, token, err)
		}
		if textMatcher {
			return definition.NewTextPattern(unescaped), nil
		}
		return definition.NewPatternAttributeMatcher(unescaped), nil
	case PegDot:
		return definition.NewDot(), nil
	case PegBuiltinSymbol:
		symbol := string(atom.SelectText())
		switch symbol {
		case definition.StartOfFileBuiltinSymbol:
			return definition.StartOfFile{}, nil
		case definition.EndOfFileBuiltinSymbol:
			return definition.EndOfFile{}, nil
		default:
			panic(fmt.Errorf("unexpected builtin symbol: %v", symbol))
		}
	}
	return nil, fmt.Errorf("can't convert atom to expression: %v", atom.Symbol)
}
