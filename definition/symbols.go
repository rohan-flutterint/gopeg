package definition

import "strings"

func AnalyzeSymbolName(symbol string) (string, bool) {
	hidden := strings.HasPrefix(symbol, "#")
	inlineDelimiter := strings.Index(symbol, "@")
	if inlineDelimiter == -1 {
		return symbol, hidden
	}
	return symbol[:inlineDelimiter], hidden
}
