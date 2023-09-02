package analysis

import (
	"fmt"
	"gopeg/definition"
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

func CheckExprConsistency(expr definition.Expr) (TerminalType, error) {
	if _, ok := expr.(definition.Terminals); !ok {
		terminalType := AnyTerminalType
		children := expr.Children()
		for i, child := range children {
			childTerminalType, err := CheckExprConsistency(child)
			if err != nil {
				return 0, err
			}
			if childTerminalType == AnyTerminalType {
				continue
			}
			if terminalType == AnyTerminalType {
				terminalType = childTerminalType
				continue
			}
			if terminalType != childTerminalType {
				return 0, fmt.Errorf("terminal type differs for expressions '%v' and '%v': %v != %v", children[i-1], children[i], terminalType, childTerminalType)
			}
		}
		return terminalType, nil
	}
	switch expr.(type) {
	case definition.Empty:
		return AnyTerminalType, nil
	case definition.Dot:
		return AnyTerminalType, nil
	case definition.TextToken:
		return ByteTerminalType, nil
	case definition.TextPattern:
		return ByteTerminalType, nil
	case definition.AtomPattern:
		return AtomTerminalType, nil
	default:
		panic(fmt.Errorf("unexpected peg expression type: %v", expr))
	}
}

func CheckRuleConsistency(rule definition.Rule) (TerminalType, error) {
	return CheckExprConsistency(rule.Expr)
}

func CheckRulesConsistency(rules definition.Rules) (TerminalType, error) {
	terminalType := AnyTerminalType
	for i, rule := range rules {
		ruleTerminalType, err := CheckRuleConsistency(rule)
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
	return terminalType, nil
}
