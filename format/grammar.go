package format

import "gopeg/parser"

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
	PegMapKeyValue = "#MapKeyValue"
	PegMapKey      = "MapKey"
	PegMapValue    = "MapValue"
)

var (
	PegGrammarRules = parser.Rules{
		parser.NewRule(PegDefinitions, parser.NewRepetition(parser.NewJunction(
			parser.NewJunction(parser.NewOptional(parser.NewSymbol(PegDefinition)), parser.NewAttr(parser.Attrs{"EndOfLine": nil})),
		))),
		parser.NewRule(PegDefinition, parser.NewJunction(
			parser.NewSymbol(PegName),
			parser.NewAttr(parser.Attrs{"Control": parser.NewToken(":")}),
			parser.NewSymbol(PegRule),
		)),
		parser.NewRule(PegName, parser.NewAttr(parser.Attrs{"Token": nil})),
		parser.NewRule(PegRule, parser.NewJunction(
			parser.NewSymbol(PegChoice),
			parser.NewRepetition(parser.NewJunction(
				parser.NewAttr(parser.Attrs{"Control": parser.NewToken("/")}),
				parser.NewSymbol(PegChoice),
			)),
		)),
		parser.NewRule(PegChoice, parser.NewRepetitionN(parser.NewSymbol(PegJunction), 1)),
		parser.NewRule(PegJunction, parser.NewJunction(
			parser.NewOptional(parser.NewJunction(parser.NewSymbol(PegAlias), parser.NewAttr(parser.Attrs{"Control": parser.NewToken(":")}))),
			parser.NewOptional(parser.NewSymbol(PegPrefix)),
			parser.NewSymbol(PegExpression),
			parser.NewOptional(parser.NewSymbol(PegSuffix)),
		)),
		parser.NewRule(PegAlias, parser.NewAttr(parser.Attrs{"Token": nil})),
		parser.NewRule(PegPrefix, parser.NewAttr(parser.Attrs{"Control": parser.NewMatch("[!&]")})),
		parser.NewRule(PegExpression, parser.NewChoice(
			parser.NewAttr(parser.Attrs{"String": nil}),
			parser.NewAttr(parser.Attrs{"Token": nil}),
			parser.NewAttr(parser.Attrs{"Regex": nil}),
			parser.NewAttr(parser.Attrs{"Any": nil}),
			parser.NewSymbol(PegMap),
			parser.NewJunction(
				parser.NewAttr(parser.Attrs{"Open": nil}),
				parser.NewSymbol(PegRule),
				parser.NewAttr(parser.Attrs{"Close": nil}),
			),
		)),
		parser.NewRule(PegSuffix, parser.NewAttr(parser.Attrs{"Control": parser.NewMatch("[+*?]")})),
		parser.NewRule(PegMap, parser.NewJunction(
			parser.NewAttr(parser.Attrs{"Control": parser.NewToken("{")}),
			parser.NewSymbol(PegMapKeyValue),
			parser.NewRepetition(parser.NewJunction(
				parser.NewAttr(parser.Attrs{"Control": parser.NewToken(",")}),
				parser.NewSymbol(PegMapKeyValue),
			)),
			parser.NewAttr(parser.Attrs{"Control": parser.NewToken("}")}),
		)),
		parser.NewRule(PegMapKeyValue, parser.NewJunction(
			parser.NewSymbol(PegMapKey),
			parser.NewOptional(parser.NewJunction(
				parser.NewAttr(parser.Attrs{"Control": parser.NewToken(":")}),
				parser.NewSymbol(PegMapValue),
			)),
		)),
		parser.NewRule(PegMapKey, parser.NewChoice(parser.NewAttr(parser.Attrs{"String": nil}), parser.NewAttr(parser.Attrs{"Token": nil}))),
		parser.NewRule(PegMapValue, parser.NewChoice(parser.NewAttr(parser.Attrs{"String": nil}), parser.NewAttr(parser.Attrs{"Regex": nil}))),
	}
)
