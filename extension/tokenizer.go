package extension

import (
	"gopeg/definition"
)

const (
	PegText      = "Text"
	PegSequence  = "#Sequence"
	PegString    = "String"
	PegRegex     = "Regex"
	PegToken     = "Token"
	PegControl   = "Control"
	PegDot       = "Dot"
	PegEndOfLine = "EndOfLine"
	PegOpen      = "Open"
	PegClose     = "Close"
)

var (
	PegTokenizerRules = definition.Rules{
		definition.NewRule(PegText, definition.NewRepetition(definition.NewJunction(
			definition.NewRepetition(definition.NewSymbol(PegSequence)),
			definition.NewSymbol(PegEndOfLine),
		))),
		definition.NewRule(PegOpen, definition.NewTextToken("(")),
		definition.NewRule(PegClose, definition.NewTextToken(")")),
		definition.NewRule(PegDot, definition.NewTextToken(".")),
		definition.NewRule(PegSequence, definition.NewRepetitionN(definition.NewChoice(
			definition.NewTextPattern("[\t\r ]+"),
			definition.NewTextPattern("//[^\n]+"),
			definition.NewJunction(
				definition.NewTextToken("/*"),
				definition.NewRepetition(definition.NewJunction(definition.NewNegation(definition.NewTextToken("*/")), definition.NewDot())),
				definition.NewTextToken("*/"),
			),
			definition.NewJunction(definition.NewTextToken("=~"), definition.NewSymbol(PegRegex)),
			definition.NewSymbol(PegString),
			definition.NewSymbol(PegToken),
			definition.NewSymbol(PegControl),
			definition.NewSymbol(PegDot),
			definition.NewJunction(
				definition.NewSymbol(PegOpen),
				definition.NewRepetition(definition.NewChoice(
					definition.NewSymbol(PegSequence),
					definition.NewTextToken("\n"),
				)),
				definition.NewSymbol(PegClose),
			),
		), 1)),
		definition.NewRule(PegRegex, definition.NewSymbol(PegString)),
		definition.NewRule(PegString, definition.NewChoice(
			definition.NewTextPattern(`'(\\.|[^'\\])*'`),
			definition.NewTextPattern(`"(\\.|[^"\\])*"`),
		)),
		definition.NewRule(PegToken, definition.NewTextPattern("[#a-zA-Z][0-9a-zA-Z_]*")),
		definition.NewRule(PegControl, definition.NewTextPattern("[:/*+?{},!&]")),
		definition.NewRule(PegEndOfLine, definition.NewChoice(
			definition.NewTextToken("\n"),
			definition.NewNegation(definition.NewDot())),
		),
	}
)
