package definition

import (
	"regexp"
	"testing"
)

func TestRegex(t *testing.T) {
	a := regexp.MustCompile("^(?i)a|ab").FindIndex([]byte("ABC"))
	t.Logf("%v", a)
}
