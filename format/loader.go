package format

import (
	"fmt"
	"gopeg/parser"
	"strconv"
)

type Rule struct {
	Choice []struct {
		Junction []struct {
			Alias      *parser.PegBindRef `peg:"Alias"`
			Prefix     *parser.PegBindRef `peg:"Prefix"`
			Suffix     *parser.PegBindRef `peg:"Suffix"`
			Expression struct {
				parser.PegBindRef
				PegRule *Rule `peg:"Rule"`
				PegMap  []struct {
					Key   *parser.PegBindRef `peg:"MapKey"`
					Value *parser.PegBindRef `peg:"MapValue"`
				} `peg:"Map"`
			} `peg:"Expression"`
		} `peg:"Junction"`
	} `peg:"Choice"`
}

type Definitions struct {
	Definitions []struct {
		Name parser.PegBindRef `peg:"Name"`
		Rule Rule              `peg:"Rule"`
	} `peg:"Definition"`
}

func Load(text string) (parser.Rules, error) {
	tokens, err := parser.Parse[byte](PegTokenizerRules, PegText, []byte(text))
	if err != nil {
		return nil, fmt.Errorf("unable to parse syntax structure: %w", err)
	}
	//fmt.Printf("%v\n", parser.StringParsingNode(tokens, []byte(text)))
	if tokens.Length() != len(text) {
		return nil, fmt.Errorf("unable to tokenize input")
	}
	attributes := make([]map[string]any, 0)
	for _, children := range tokens.Children() {
		attributes = append(attributes, children.Attrs())
	}
	//fmt.Printf("attrs: %v\n", attributes)
	peg, err := parser.Parse[map[string]any](PegGrammarRules, PegDefinitions, attributes)
	//fmt.Printf("peg: %v\n", parser.StringParsingNode(peg, attributes))
	if err != nil {
		return nil, fmt.Errorf("unable to parse lexical structure: %w", err)
	}
	if peg.Length() != len(attributes) {
		return nil, fmt.Errorf("unable to parse input")
	}
	var ds Definitions
	err = parser.Bind(peg, &ds)
	if err != nil {
		return nil, err
	}
	rules := make(parser.Rules, 0)
	for _, d := range ds.Definitions {
		name, ok := token(attributes[d.Name.Start][PegToken])
		if !ok {
			return nil, fmt.Errorf("unexpected name")
		}
		current, additional, err := rule(d.Rule, attributes)
		if err != nil {
			return nil, err
		}
		rules = append(rules, additional...)
		rules = append(rules, parser.NewRule(name, current))
	}
	return rules, nil
}

func rule(r Rule, attributes []map[string]any) (parser.Expr, []parser.Rule, error) {
	rules := make([]parser.Rule, 0)
	choices := make([]parser.Expr, 0, len(r.Choice))
	for _, c := range r.Choice {
		junctions := make([]parser.Expr, 0, len(c.Junction))
		for _, j := range c.Junction {
			var current parser.Expr
			var err error
			if len(j.Expression.PegMap) > 0 {
				attrs := make(parser.Attrs)
				for _, entry := range j.Expression.PegMap {
					var key string
					if keyString, ok := str(attributes[entry.Key.Start]["String"]); ok {
						key = keyString
					}
					if keyToken, ok := token(attributes[entry.Key.Start]["Token"]); ok {
						key = keyToken
					}
					if key == "" {
						return nil, nil, fmt.Errorf("empty key")
					}
					if entry.Value != nil {
						value, err := expr(attributes[entry.Value.Start])
						if err == nil {
							attrs[key] = value
						} else {
							return nil, nil, err
						}
					} else {
						attrs[key] = nil
					}
				}
				current = parser.NewAttr(attrs)
			} else if j.Expression.PegRule != nil {
				var addition []parser.Rule
				current, addition, err = rule(*j.Expression.PegRule, attributes)
				if err != nil {
					return nil, nil, err
				}
				rules = append(rules, addition...)
			} else {
				current, err = expr(attributes[j.Expression.Start])
				if err != nil {
					return nil, nil, err
				}
			}
			if j.Prefix != nil {
				control, ok := token(attributes[j.Prefix.Start]["Control"])
				if !ok {
					return nil, nil, fmt.Errorf("unexpected prefix")
				}
				if control == "!" {
					current = parser.NewNegation(current)
				} else if control == "&" {
					current = parser.NewEnsure(current)
				} else {
					return nil, nil, fmt.Errorf("unknown prefix: %v", control)
				}
			}
			if j.Suffix != nil {
				control, ok := token(attributes[j.Suffix.Start]["Control"])
				if !ok {
					return nil, nil, fmt.Errorf("unexpected suffix")
				}
				if control == "?" {
					current = parser.NewOptional(current)
				} else if control == "*" {
					current = parser.NewRepetition(current)
				} else if control == "+" {
					current = parser.NewRepetitionN(current, 1)
				} else {
					return nil, nil, fmt.Errorf("unknown suffix: %v", control)
				}
			}
			if j.Alias != nil {
				alias, ok := token(attributes[j.Alias.Start]["Token"])
				if !ok {
					return nil, nil, fmt.Errorf("unexpected alias")
				}
				rules = append(rules, parser.NewRule(alias, current))
				current = parser.NewSymbol(alias)
			}
			junctions = append(junctions, current)
		}
		choices = append(choices, parser.NewJunction(junctions...))
	}
	return parser.NewChoice(choices...), rules, nil
}

func token(s any) (string, bool) {
	if _, ok := s.([]byte); !ok {
		return "", false
	}
	return string(s.([]byte)), true
}

func str(s any) (string, bool) {
	if _, ok := s.([]byte); !ok {
		return "", false
	}
	str, err := unescape(string(s.([]byte)))
	if err != nil {
		return "", false
	}
	return str, true
}

func expr(e map[string]any) (parser.Expr, error) {
	if token, ok := token(e["Token"]); ok {
		return parser.NewSymbol(token), nil
	} else if s, ok := str(e["String"]); ok {
		return parser.NewToken(s), nil
	} else if regex, ok := str(e["Regex"]); ok {
		return parser.NewMatch(regex), nil
	} else if _, ok := e["Any"]; ok {
		return parser.NewAny(), nil
	}
	return nil, fmt.Errorf("unknown expression")
}

func unescape(s string) (string, error) {
	return strconv.Unquote(s)
}
