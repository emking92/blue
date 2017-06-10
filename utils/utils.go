package utils

import (
	"fmt"
	"regexp"
)

type StrSubstititor map[string]func(string) string

var (
	tokenRegex *regexp.Regexp
)

func init() {
	tokenRegex, _ = regexp.Compile(`\b([A-Za-z_][A-Za-z0-9_]*)\b`)
}

func (vr StrSubstititor) CreateVariable(v string, value string) {
	reg := fmt.Sprintf(`\b%s\b`, v)

	varRegex, err := regexp.Compile(reg)
	if err != nil {
		panic(err)
	}

	vr[v] = func(str string) string {
		return varRegex.ReplaceAllLiteralString(str, value)
	}
}

func (vr StrSubstititor) Expand(str string) (out string) {
	out = str

	matches := tokenRegex.FindAllStringSubmatch(str, -1)

	for _, m := range matches {
		token := m[1]
		expander, ok := vr[token]
		if !ok {
			continue
		}

		out = expander(out)
	}

	return
}

func (vr StrSubstititor) IsVarDefined(str string) bool {
	_, ok := vr[str]
	return ok
}
