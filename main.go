package main

import (
	"fmt"
	"gopeg/format"
	"gopeg/parser"
	"os"
)

func main() {
	tokenizerFile, err := os.ReadFile("sql-tokenizer.peg")
	if err != nil {
		panic(err)
	}
	tokenizerRules, err := format.Load(string(tokenizerFile))
	if err != nil {
		panic(err)
	}
	sql, err := os.ReadFile("test.sql")
	if err != nil {
		panic(err)
	}
	root, err := parser.Parse(tokenizerRules, "Text", sql)
	if err != nil {
		panic(err)
	}
	attributes := make([]map[string]any, 0)
	for _, children := range root.Children() {
		attributes = append(attributes, children.Attrs())
	}
	//fmt.Printf("attrs: %v\n", attributes)
	grammarFile, err := os.ReadFile("sql-grammar.peg")
	if err != nil {
		panic(err)
	}
	grammarRules, err := format.Load(string(grammarFile))
	if err != nil {
		panic(err)
	}
	//fmt.Printf("grammar: %v\n", grammarRules)
	peg, err := parser.Parse[map[string]any](grammarRules, "Expressions", attributes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", parser.StringParsingNode(peg, attributes))
}
