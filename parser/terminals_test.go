package parser

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

func TestFormat(t *testing.T) {
	assert.Equal(t, "\"a\\n\\r\\tb\"", format([]byte("a\n\r\tb")))
	assert.Equal(t, "'\\n'", format(byte('\n')))
}

func TestMatcherFormat(t *testing.T) {
	assert.Equal(t, "=~\"[^\nabc]\"", byteMatcher{expr: "[^\nabc]", regex: regexp.MustCompile("[^\nabc]")}.String())
}

func TestDotFormat(t *testing.T) {
	assert.Equal(t, ".", dot{}.String())
}

func TestSequenceFormat(t *testing.T) {
	assert.Equal(t, "\"hi\"", byteSequence{value: []byte("hi")}.String())
}

func TestAttrFormat(t *testing.T) {
	assert.Equal(t, "[{Control, Attr:\"Test\"}]", attrSequenceMatcher([]attrMatcher{
		{"Control": nil, "Attr": byteSequence{value: []byte("Test")}},
	}).String())
}

func TestRegex(t *testing.T) {
	a := regexp.MustCompile("^(?i)a|ab").FindIndex([]byte("ABC"))
	t.Logf("%v", a)
}

func BenchmarkTokenAccept(b *testing.B) {
	t := byteSequence{value: []byte(strings.Repeat("a", 128))}
	text := make([]byte, 0)
	for i := 0; i < 1024; i++ {
		text = append(text, []byte(strings.Repeat("a", 1024))...)
		text = append(text, 'b')
	}
	cnt := 0
	for i := 0; i < b.N; i++ {
		s, ok := t.Accept(text[i%2048:])
		if ok {
			cnt += s
		}
	}
	if cnt == 0 {
		b.Fail()
	}
}
