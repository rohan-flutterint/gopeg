package format

import "gopeg/parser"

const (
	PegText      = "Text"
	PegSequence  = "#Sequence"
	PegString    = "String"
	PegRegex     = "Regex"
	PegToken     = "Token"
	PegControl   = "Control"
	PegAny       = "Any"
	PegEndOfLine = "EndOfLine"
	PegOpen      = "Open"
	PegClose     = "Close"
)

var (
	PegTokenizerRules = parser.Rules{
		parser.NewRule(PegText, parser.NewRepetition(parser.NewJunction(
			parser.NewRepetition(parser.NewSymbol(PegSequence)),
			parser.NewSymbol(PegEndOfLine),
		))),
		parser.NewRule(PegOpen, parser.NewToken("(")),
		parser.NewRule(PegClose, parser.NewToken(")")),
		parser.NewRule(PegAny, parser.NewToken(".")),
		parser.NewRule(PegSequence, parser.NewRepetitionN(parser.NewChoice(
			parser.NewMatch("[\t\r ]+"),
			parser.NewMatch("//[^\n]+"),
			parser.NewJunction(
				parser.NewToken("/*"),
				parser.NewRepetition(parser.NewJunction(parser.NewNegation(parser.NewToken("*/")), parser.NewAny())),
				parser.NewToken("*/"),
			),
			parser.NewJunction(parser.NewToken("=~"), parser.NewSymbol(PegRegex)),
			parser.NewSymbol(PegString),
			parser.NewSymbol(PegToken),
			parser.NewSymbol(PegControl),
			parser.NewSymbol(PegAny),
			parser.NewJunction(
				parser.NewSymbol(PegOpen),
				parser.NewRepetition(parser.NewChoice(
					parser.NewSymbol(PegSequence),
					parser.NewToken("\n"),
				)),
				parser.NewSymbol(PegClose),
			),
		), 1)),
		parser.NewRule(PegRegex, parser.NewSymbol(PegString)),
		parser.NewRule(PegString, parser.NewChoice(
			parser.NewMatch(`'(\\.|[^'\\])*'`),
			parser.NewMatch(`"(\\.|[^"\\])*"`),
		)),
		parser.NewRule(PegToken, parser.NewMatch("[#a-zA-Z][0-9a-zA-Z_]*")),
		parser.NewRule(PegControl, parser.NewMatch("[:/*+?{},!&]")),
		parser.NewRule(PegEndOfLine, parser.NewChoice(
			parser.NewToken("\n"),
			parser.NewNegation(parser.NewAny())),
		),
	}
)
