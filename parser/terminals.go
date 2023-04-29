package parser

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type terminals interface {
	fmt.Stringer
	Accept(seq any) (int, bool)
}

type (
	dot          struct{}
	empty        struct{}
	byteSequence struct{ value []byte }
	byteMatcher  struct {
		regex *regexp.Regexp
		expr  string
	}

	attrMatcher         map[string]terminals
	attrSequenceMatcher []attrMatcher
)

func (a attrMatcher) String() string {
	entries := make([]string, 0)
	nonNilKeys := make([]string, 0)
	for key, value := range a {
		if value == nil {
			entries = append(entries, key)
		} else {
			nonNilKeys = append(nonNilKeys, key)
		}
	}
	sort.Strings(entries)
	sort.Strings(nonNilKeys)
	for _, key := range nonNilKeys {
		entries = append(entries, fmt.Sprintf("%v:%v", key, a[key]))
	}
	return fmt.Sprintf("{%v}", strings.Join(entries, ", "))
}

func prefixes[T comparable](a, b any) (int, bool) {
	aBytes, aOk := a.([]byte)
	bBytes, bOk := b.([]byte)
	if aOk && bOk {
		if len(aBytes) <= len(bBytes) && string(aBytes) == string(bBytes[:len(aBytes)]) {
			return len(aBytes), true
		}
		return 0, false
	}
	aT := a.([]T)
	bT := b.([]T)
	if len(aT) > len(bT) {
		return 0, false
	}
	for i := 0; i < len(aT); i++ {
		if aT[i] != bT[i] {
			return 0, false
		}
	}
	return len(aT), true
}

func format(seq any) string {
	t := reflect.TypeOf(seq)
	if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		quoted := strconv.Quote(string(seq.([]byte)))
		return quoted
	}
	if t.Kind() == reflect.Uint8 {
		quoted := strconv.QuoteRune(rune(seq.(byte)))
		return quoted
	}
	return fmt.Sprintf("%v", seq)
}

func (a attrSequenceMatcher) Accept(seq any) (int, bool) {
	maps := seq.([]map[string]any)
	if len(maps) < len(a) {
		return 0, false
	}
	//fmt.Printf("try match: %v against %v\n", maps, a)
	for i := range a {
		for key, matcher := range a[i] {
			value, ok := maps[i][key]
			if !ok {
				return 0, false
			}
			if matcher == nil {
				continue
			}
			if _, ok = matcher.Accept(value); !ok {
				return 0, false
			}
		}
	}
	return len(a), true
}

func (a attrSequenceMatcher) String() string { return fmt.Sprintf("%v", []attrMatcher(a)) }

func (t byteSequence) Accept(seq any) (int, bool) {
	return prefixes[byte](t.value, seq)
}
func (t byteSequence) String() string { return format(t.value) }

func (m byteMatcher) Accept(s any) (int, bool) {
	a := s.([]byte)
	location := m.regex.FindIndex(a)
	if location == nil || location[0] != 0 {
		return 0, false
	}
	return location[1], true
}
func (m byteMatcher) String() string { return fmt.Sprintf("=~\"%v\"", m.expr) }

func (d dot) Accept(s any) (int, bool) {
	if reflect.ValueOf(s).Len() > 0 {
		return 1, true
	}
	return 0, false
}
func (d dot) String() string { return "." }

func (e empty) Accept(_ any) (int, bool) { return 0, true }
func (e empty) String() string           { return "e" }
