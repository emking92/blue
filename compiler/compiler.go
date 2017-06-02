package compiler

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"
)

var (
	lineRegex           *regexp.Regexp
	lineRegexIndexLabel int
	lineRegexIndexOp    int
	lineRegexIndexArg1  int
	lineRegexIndexArg2  int
	lineRegexIndexArg3  int

	ignoredLineRegex *regexp.Regexp
)

func init() {
	var err error
	lineRegex, err = regexp.Compile(`^\s*(\w[\w\d+]*\s*:)?\s*([A-Za-z]+)\s+([\w\d\[\]\$]+)(\s*,\s*([\w\d\[\]\$]+)(\s*,\s*([\w\d\[\]\$]+))?)?\s*(;.*)?$`)
	lineRegexIndexLabel = 1
	lineRegexIndexOp = 2
	lineRegexIndexArg1 = 3
	lineRegexIndexArg2 = 5
	lineRegexIndexArg3 = 7

	if err != nil {
		panic(err)
	}

	ignoredLineRegex, err = regexp.Compile(`^\s*(;.*)?$`)
}

type codeParts struct {
	label string
	op    string
	args  []string
}

func BuildSource(source io.Reader) (instructions []Instruction, err error) {
	scanner := bufio.NewScanner(source)
	if err != nil {
		return nil, err
	}

	pgm := programBuilder{}

	//	syntaxErrors := []error{}

	for scanner.Scan() {
		line := scanner.Text()
		pgm.currentLineNumber++

		//Ignore empty lines and lines with only comments
		if ignoredLineRegex.MatchString(line) {
			continue
		}

		pgm.currentInstructionNumber++

		matches := lineRegex.FindStringSubmatch(line)
		if matches == nil {
			pgm.compileErrorString("syntax error: " + strings.TrimSpace(line))
			continue
		}

		parts := codeParts{
			label: matches[lineRegexIndexLabel],
			op:    matches[lineRegexIndexOp],
			args: []string{
				matches[lineRegexIndexArg1],
				matches[lineRegexIndexArg2],
				matches[lineRegexIndexArg3],
			},
		}

		pgm.buildInstruction(parts)

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if pgm.errs != nil {
		errs := append([]string{"Compilation errors: "}, pgm.errs...)
		err = errors.New(strings.Join(errs, "\n\t"))
	}

	return pgm.instructions, err
}
