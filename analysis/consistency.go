package analysis

import (
	"fmt"
	"github.com/sivukhin/gopeg/definition"
)

type TerminalType int

const (
	AnyTerminalType  TerminalType = 0
	ByteTerminalType TerminalType = 1
	AtomTerminalType TerminalType = 2
)

func (t TerminalType) String() string {
	switch t {
	case AnyTerminalType:
		return "AnyTerminalType"
	case ByteTerminalType:
		return "ByteTerminalType"
	case AtomTerminalType:
		return "AtomTerminalType"
	default:
		panic(fmt.Errorf("unexpected terminal type: %v", int(t)))
	}
}

func CheckExprConsistency(expr definition.Expr) (TerminalType, map[string]struct{}, error) {
	if _, ok := expr.(definition.Terminals); !ok {
		terminalType := AnyTerminalType
		children := expr.Children()
		ruleNames := make(map[string]struct{})
		if symbol, ok := expr.(definition.Symbol); ok {
			ruleNames[symbol.Name] = struct{}{}
		}
		for i, child := range children {
			childTerminalType, childRuleNames, err := CheckExprConsistency(child)
			for ruleName := range childRuleNames {
				ruleNames[ruleName] = struct{}{}
			}
			if err != nil {
				return 0, nil, err
			}
			if childTerminalType == AnyTerminalType {
				continue
			}
			if terminalType == AnyTerminalType {
				terminalType = childTerminalType
				continue
			}
			if terminalType != childTerminalType {
				return 0, nil, fmt.Errorf("terminal type differs for expressions '%v' and '%v': %v != %v", children[i-1], children[i], terminalType, childTerminalType)
			}
		}
		return terminalType, ruleNames, nil
	}
	switch expr.(type) {
	case definition.Empty:
		return AnyTerminalType, nil, nil
	case definition.Dot:
		return AnyTerminalType, nil, nil
	case definition.TextToken:
		return ByteTerminalType, nil, nil
	case definition.TextPattern:
		return ByteTerminalType, nil, nil
	case definition.AtomPattern:
		return AtomTerminalType, nil, nil
	case definition.StartOfFile:
		return AnyTerminalType, nil, nil
	case definition.EndOfFile:
		return AnyTerminalType, nil, nil
	default:
		panic(fmt.Errorf("unexpected peg expression type: %v", expr))
	}
}

func CheckRuleConsistency(rule definition.Rule) (TerminalType, map[string]struct{}, error) {
	return CheckExprConsistency(rule.Expr)
}

func CheckRulesConsistency(rules definition.Rules) (TerminalType, error) {
	terminalType := AnyTerminalType
	usedRuleNames := make(map[string]struct{})
	definedRuleNames := make(map[string]struct{})
	for i, rule := range rules {
		definedRuleNames[rule.Name] = struct{}{}
		ruleTerminalType, ruleNames, err := CheckRuleConsistency(rule)
		for ruleName := range ruleNames {
			usedRuleNames[ruleName] = struct{}{}
		}

		if err != nil {
			return 0, fmt.Errorf("inconsistent terminal type in rule '%v': %w", rule.Name, err)
		}
		if ruleTerminalType == AnyTerminalType {
			continue
		}
		if terminalType == AnyTerminalType {
			terminalType = ruleTerminalType
			continue
		}
		if terminalType != ruleTerminalType {
			return 0, fmt.Errorf("terminal type differs for rules '%v' and '%v': %v != %v", rules[i-1], rules[i], terminalType, ruleTerminalType)
		}
	}
	for ruleName := range usedRuleNames {
		if _, ok := definedRuleNames[ruleName]; !ok {
			return terminalType, fmt.Errorf("undefined rule '%v' were used in the expression", ruleName)
		}
	}
	return terminalType, nil
}
