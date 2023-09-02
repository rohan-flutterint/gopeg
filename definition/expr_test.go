package definition

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAtom(t *testing.T) {
	atom := Atom{
		Attributes: nil,
		TextSelector: BuildSegments(
			Segment{Start: 3, End: 5},
			Segment{Start: 9, End: 10},
			Segment{Start: 11, End: 12},
		),
		Text: []byte("hello, world"),
	}
	require.Equal(t, []byte("lord"), atom.SelectText())
}

func TestExprString(t *testing.T) {
	t.Run("precedence check", func(t *testing.T) {
		require.Equal(t, `("a" / "b") ("c" / "d")`, NewJunction(
			NewChoice(NewTextToken("a"), NewTextToken("b")),
			NewChoice(NewTextToken("c"), NewTextToken("d")),
		).String())
	})
	t.Run("all node types", func(t *testing.T) {
		require.Equal(t, `@empty / . / "a" / {Text:=~"^[0-9]+$"} / {Text:"test"} / =~"^[0-9]*" / ("a" "b")* / ("a" "b")+ / ("a" "b"){2,} / ("a" "b")? / &"a"* / (&"a")* / !"a"* / (!"a")* / A`, NewChoice(
			NewEmpty(),
			NewDot(),
			NewTextToken("a"),
			NewAtomPattern(map[string]AttributeMatcher{"Text": NewPatternAttributeMatcher("[0-9]+")}),
			NewAtomPattern(map[string]AttributeMatcher{"Text": NewTokenAttributeMatcher("test")}),
			NewTextPattern("[0-9]*"),
			NewRepetition(NewJunction(NewTextToken("a"), NewTextToken("b"))),
			NewRepetitionN(NewJunction(NewTextToken("a"), NewTextToken("b")), 1),
			NewRepetitionN(NewJunction(NewTextToken("a"), NewTextToken("b")), 2),
			NewOptional(NewJunction(NewTextToken("a"), NewTextToken("b"))),
			NewEnsure(NewRepetition(NewTextToken("a"))),
			NewRepetition(NewEnsure(NewTextToken("a"))),
			NewNegation(NewRepetition(NewTextToken("a"))),
			NewRepetition(NewNegation(NewTextToken("a"))),
			NewSymbol("A"),
		).String())
	})
}
