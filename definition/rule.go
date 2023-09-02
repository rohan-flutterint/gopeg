package definition

import (
	"fmt"
	"strings"
)

type (
	Rule struct {
		Name string
		Expr Expr
	}
	Rules []Rule
)

func (rs Rules) Combine(other ...Rules) Rules {
	combined := append(Rules{}, rs...)
	for _, b := range other {
		combined = append(rs, b...)
	}
	return combined
}
func NewRule(name string, expression Expr) Rule {
	return Rule{Name: name, Expr: expression}
}
func (r Rule) String() string {
	return fmt.Sprintf("%v: %v", r.Name, r.Expr)
}

func (rs Rules) String() string {
	var b strings.Builder
	for _, r := range rs {
		b.WriteString(r.String())
		b.WriteString("\n")
	}
	return b.String()
}
