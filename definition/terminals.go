package definition

import (
	"bytes"
	"fmt"
)

func (e TextToken) Match(data []byte) (int, bool) {
	if bytes.HasPrefix(data, e.Text) {
		return len(e.Text), true
	}
	return 0, false
}
func (e TextPattern) Match(data []byte) (int, bool) {
	location := e.Regex.FindIndex(data)
	if location == nil || location[0] != 0 {
		return 0, false
	}
	return location[1], true
}

func Accept[T any](terminal Terminals, text []T) (int, bool) {
	switch peg := terminal.(type) {
	case Empty:
		return 0, true
	case Dot:
		if len(text) == 0 {
			return 0, false
		}
		return 1, true
	case TextTerminals:
		textBytes, textOk := any(text).([]byte)
		if !textOk {
			panic(fmt.Errorf("TextToken terminal can be used only for byte sequences, given %#v", terminal))
		}
		return peg.Match(textBytes)
	case AtomPattern:
		atoms, atomsOk := any(text).([]Atom)
		if !atomsOk {
			panic(fmt.Errorf("AtomPattern terminal can be used only for atom sequences, given %T", *new(T)))
		}
		if len(atoms) == 0 {
			return 0, false
		}
		textAtom := atoms[0]
		for key, matcher := range peg.Matcher {
			if key == textAtom.Symbol {
				if textAtom.Symbol != key {
					return 0, false
				}
				if matcher == nil {
					continue
				}
				if _, ok := matcher.Match(textAtom.SelectText()); !ok {
					return 0, false
				}
				continue
			}
			attribute, exists := textAtom.Attributes[key]
			if !exists {
				return 0, false
			}
			if matcher == nil {
				continue
			}
			if _, ok := matcher.Match(attribute); ok {
				return 0, false
			}
		}
		return 1, true
	default:
		panic(fmt.Errorf("unexpected peg expression type: %#v", terminal))
	}
}
