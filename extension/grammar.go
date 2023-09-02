package extension

import (
	"gopeg/definition"
)

const (
	PegDefinitions = "Definitions"
	PegDefinition  = "Definition"
	PegName        = "Name"
	PegRule        = "Rule"
	PegChoice      = "Choice"
	PegJunction    = "Junction"
	PegAlias       = "Alias"
	PegPrefix      = "Prefix"
	PegExpression  = "Expression"
	PegSuffix      = "Suffix"
	PegMap         = "Map"
	PegMapKeyValue = "MapKeyValue"
	PegMapKey      = "MapKey"
	PegMapValue    = "MapValue"
)

var (
	PegGrammarRules = definition.Rules{
		definition.NewRule(PegDefinitions, definition.NewRepetition(definition.NewJunction(
			definition.NewOptional(definition.NewSymbol(PegDefinition)),
			definition.NewAtomPattern(map[string]definition.TextTerminals{PegEndOfLine: nil}),
		))),
		definition.NewRule(PegDefinition, definition.NewJunction(
			definition.NewSymbol(PegName),
			definition.NewAtomPattern(map[string]definition.TextTerminals{PegControl: definition.NewTokenAttributeMatcher(":")}),
			definition.NewSymbol(PegRule),
		)),
		definition.NewRule(PegName, definition.NewAtomPattern(map[string]definition.TextTerminals{PegToken: nil})),
		definition.NewRule(PegRule, definition.NewJunction(
			definition.NewSymbol(PegChoice),
			definition.NewRepetition(definition.NewJunction(
				definition.NewAtomPattern(map[string]definition.TextTerminals{PegControl: definition.NewTokenAttributeMatcher("/")}),
				definition.NewSymbol(PegChoice),
			)),
		)),
		definition.NewRule(PegChoice, definition.NewRepetitionN(definition.NewSymbol(PegJunction), 1)),
		definition.NewRule(PegJunction, definition.NewJunction(
			definition.NewOptional(definition.NewJunction(definition.NewSymbol(PegAlias), definition.NewAtomPattern(map[string]definition.TextTerminals{PegControl: definition.NewTokenAttributeMatcher(":")}))),
			definition.NewOptional(definition.NewSymbol(PegPrefix)),
			definition.NewSymbol(PegExpression),
			definition.NewOptional(definition.NewSymbol(PegSuffix)),
		)),
		definition.NewRule(PegAlias, definition.NewAtomPattern(map[string]definition.TextTerminals{PegToken: nil})),
		definition.NewRule(PegPrefix, definition.NewAtomPattern(map[string]definition.TextTerminals{PegControl: definition.NewPatternAttributeMatcher("[!&]")})),
		definition.NewRule(PegExpression, definition.NewChoice(
			definition.NewAtomPattern(map[string]definition.TextTerminals{PegString: nil}),
			definition.NewAtomPattern(map[string]definition.TextTerminals{PegToken: nil}),
			definition.NewAtomPattern(map[string]definition.TextTerminals{PegRegex: nil}),
			definition.NewAtomPattern(map[string]definition.TextTerminals{PegDot: nil}),
			definition.NewSymbol(PegMap),
			definition.NewJunction(
				definition.NewAtomPattern(map[string]definition.TextTerminals{PegOpen: nil}),
				definition.NewSymbol(PegRule),
				definition.NewAtomPattern(map[string]definition.TextTerminals{PegClose: nil}),
			),
		)),
		definition.NewRule(PegSuffix, definition.NewAtomPattern(map[string]definition.TextTerminals{PegControl: definition.NewPatternAttributeMatcher("[+*?]")})),
		definition.NewRule(PegMap, definition.NewJunction(
			definition.NewAtomPattern(map[string]definition.TextTerminals{PegControl: definition.NewTokenAttributeMatcher("{")}),
			definition.NewSymbol(PegMapKeyValue),
			definition.NewRepetition(definition.NewJunction(
				definition.NewAtomPattern(map[string]definition.TextTerminals{PegControl: definition.NewTokenAttributeMatcher(",")}),
				definition.NewSymbol(PegMapKeyValue),
			)),
			definition.NewAtomPattern(map[string]definition.TextTerminals{PegControl: definition.NewTokenAttributeMatcher("}")}),
		)),
		definition.NewRule(PegMapKeyValue, definition.NewJunction(
			definition.NewSymbol(PegMapKey),
			definition.NewOptional(definition.NewJunction(
				definition.NewAtomPattern(map[string]definition.TextTerminals{PegControl: definition.NewTokenAttributeMatcher(":")}),
				definition.NewSymbol(PegMapValue),
			)),
		)),
		definition.NewRule(PegMapKey, definition.NewChoice(definition.NewAtomPattern(map[string]definition.TextTerminals{PegString: nil}), definition.NewAtomPattern(map[string]definition.TextTerminals{PegToken: nil}))),
		definition.NewRule(PegMapValue, definition.NewChoice(definition.NewAtomPattern(map[string]definition.TextTerminals{PegString: nil}), definition.NewAtomPattern(map[string]definition.TextTerminals{PegRegex: nil}))),
	}
)
