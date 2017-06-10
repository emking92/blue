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

	preprocessLineRegex          *regexp.Regexp
	preprocesslineRegexIndexOp   int
	preprocesslineRegexIndexArg1 int
	preprocesslineRegexIndexArg2 int
	preprocesslineRegexIndexArg3 int

	ignoredLineRegex *regexp.Regexp
)

func init() { //                          (label        ):     (command  )   (arg1          )    ,   (arg2          )    ,   (arg3          )        ;comments
	lineRegex, _ = regexp.Compile(`^\s*(?:(\w[\w\d+]*\s*):)?\s*([A-Za-z]+)\s+(-?[\w\d\[\]&]+)(\s*,\s*(-?[\w\d\[\]&]+)(\s*,\s*(-?[\w\d\[\]&]+))?)?\s*(;.*)?$`)
	lineRegexIndexLabel = 1
	lineRegexIndexOp = 2
	lineRegexIndexArg1 = 3
	lineRegexIndexArg2 = 5
	lineRegexIndexArg3 = 7

	preprocessLineRegex, _ = regexp.Compile(`^\s*#([A-Za-z]+)(\s+(-?[\w\d]+)(\s*,\s*(-?[\w\d]+)(\s*,\s*(-?[\w\d]+))?)?)?\s*(;.*)?$`)
	preprocesslineRegexIndexOp = 1
	preprocesslineRegexIndexArg1 = 3
	preprocesslineRegexIndexArg2 = 5
	preprocesslineRegexIndexArg3 = 7

	ignoredLineRegex, _ = regexp.Compile(`^\s*(;.*)?$`)
}

type codeParts struct {
	lineNumber int
	label      string
	op         string
	args       []string
}

type preprocessorParts struct {
	lineNumber int
	op         string
	args       []string
}

func BuildSource(source io.Reader) (instructions []Instruction, err error) {
	scanner := bufio.NewScanner(source)
	if err != nil {
		return nil, err
	}

	pgm := programBuilder{}
	pgm.init()

	var parsedCode []codeParts
	lineNumber := 0
	instructionIndex := -1

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		//Ignore empty lines and lines with only comments
		if ignoredLineRegex.MatchString(line) {
			continue
		}

		matches := preprocessLineRegex.FindStringSubmatch(line)
		if matches != nil {
			parts := preprocessorParts{
				lineNumber: lineNumber,
				op:         matches[preprocesslineRegexIndexOp],
				args: []string{
					matches[preprocesslineRegexIndexArg1],
					matches[preprocesslineRegexIndexArg2],
					matches[preprocesslineRegexIndexArg3],
				},
			}

			pgm.preprocess(parts)
			continue
		}

		instructionIndex++

		matches = lineRegex.FindStringSubmatch(line)
		if matches == nil {
			pgm.compileErrorString("syntax error: " + strings.TrimSpace(line))
			continue
		}

		parts := codeParts{
			lineNumber: lineNumber,
			label:      matches[lineRegexIndexLabel],
			op:         matches[lineRegexIndexOp],
			args: []string{
				matches[lineRegexIndexArg1],
				matches[lineRegexIndexArg2],
				matches[lineRegexIndexArg3],
			},
		}
		parsedCode = append(parsedCode, parts)

		if len(parts.label) > 0 {
			if _, ok := pgm.labels[strings.ToLower(parts.label)]; ok {
				pgm.compileErrorString(`label already defined "%s"`, parts.label)
				continue
			}
			pgm.setLabelLine(parts.label, instructionIndex)
		}

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for _, codeLines := range parsedCode {
		pgm.buildInstruction(codeLines)
	}

	if pgm.errs != nil {
		errs := append([]string{"Compilation errors: "}, pgm.errs...)
		err = errors.New(strings.Join(errs, "\n\t"))
	}

	return pgm.instructions, err
}
