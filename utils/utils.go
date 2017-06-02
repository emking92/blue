package utils

import (
	"fmt"
	"regexp"
	"strings"
)

type StrSubstititor map[string]func(string) string

var (
	varRegex *regexp.Regexp
)

func init() {
	varRegex, _ = regexp.Compile(`(?:^|\W)\$([A-Za-z_][A-Za-z0-9_]*)(?:\b|$)`)
}

func (vr StrSubstititor) CreateVariable(v string, value string) {
	reg := fmt.Sprintf(`\$%s\b`, v)

	varRegex, err := regexp.Compile(reg)
	if err != nil {
		panic(err)
	}

	vr[v] = func(str string) string {
		return varRegex.ReplaceAllLiteralString(str, value)
	}
}

func (vr StrSubstititor) Expand(str string) (out string, err error) {
	out = str

	if !strings.ContainsRune(str, '$') {
		return
	}

	matches := varRegex.FindAllStringSubmatch(str, -1)
	if len(matches) == 0 {
		err = fmt.Errorf(`invalid substitution "%s"`, str)
		return
	}

	for _, m := range matches {
		varFound := m[1]
		expander, ok := vr[varFound]

		if !ok {
			err = fmt.Errorf(`undefined var "%s"`, varFound)
			return
		}

		out = expander(out)
	}

	return
}

func (vr StrSubstititor) IsVarDefined(str string) bool {
	_, ok := vr[str]
	return ok
}
