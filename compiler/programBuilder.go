package compiler

import (
	"blue/utils"
	"fmt"
	"strings"
)

type programBuilder struct {
	instructions      []Instruction
	labels            map[string]int
	variables         utils.StrSubstititor
	nextWordAddress   int
	currentLineNumber int
	errs              []string
}

func (pgm *programBuilder) init() {
	pgm.labels = make(map[string]int)
	pgm.variables = make(utils.StrSubstititor)
	pgm.nextWordAddress = 16
}

func (pgm *programBuilder) setLabelLine(label string, line int) {
	pgm.labels[strings.ToLower(label)] = line
}

func (pgm programBuilder) getLabelLine(label string) (line int, ok bool) {
	line, ok = pgm.labels[strings.ToLower(label)]
	return
}

func (pgm programBuilder) isLabelDefined(label string) bool {
	_, ok := pgm.labels[strings.ToLower(label)]
	return ok
}

func (pgm *programBuilder) allocateWord() int {
	next := pgm.nextWordAddress
	pgm.nextWordAddress++
	return next
}

func (pgm *programBuilder) compileError(err error) {
	pgm.compileErrorString(err.Error())
}

func (pgm *programBuilder) compileErrorString(message string, args ...interface{}) {
	str1 := fmt.Sprintf("Error: %d: ", pgm.currentLineNumber)
	str := fmt.Sprintf(str1+message, args...)
	pgm.errs = append(pgm.errs, str)
}

func (pgm *programBuilder) preprocess(parts preprocessorParts) {
	pgm.currentLineNumber = parts.lineNumber
	op := strings.ToLower(parts.op)

	strategy, strategyOk := precompilerStrategies[op]
	if !strategyOk {
		pgm.compileErrorString("undefined preprocess operation: " + parts.op)
		return
	}

	argCount := 0
	if len(parts.args[2]) > 0 {
		argCount = 3
	} else if len(parts.args[1]) > 0 {
		argCount = 2
	} else if len(parts.args[0]) > 0 {
		argCount = 1
	}

	if argCount != strategy.argCount {
		pgm.compileErrorString("invalid argument count for preprocess operation %s. Requires %d, but received %v", op, strategy.argCount, parts.args[:argCount])
		return
	}

	args := make([]string, argCount)

	for i := 0; i < argCount; i++ {
		args[i] = pgm.variables.Expand(parts.args[i])
	}

	strategy.precompile(pgm, op, args)
}

func (pgm *programBuilder) buildInstruction(parts codeParts) {
	pgm.currentLineNumber = parts.lineNumber
	op := strings.ToLower(parts.op)

	strategy, strategyOk := opBuildStrategies[op]
	if !strategyOk {
		pgm.compileErrorString("undefined operation: " + parts.op)
	}

	argCount := 0
	if len(parts.args[2]) > 0 {
		argCount = 3
	} else if len(parts.args[1]) > 0 {
		argCount = 2
	} else if len(parts.args[0]) > 0 {
		argCount = 1
	}

	args := make(argumentGroup, argCount)

	var err error
	for i := 0; i < argCount; i++ {
		args[i], err = pgm.parseInstructionArgument(parts.args[i], i)
		if err != nil {
			pgm.compileErrorString(err.Error())
			continue
		}
	}

	if !strategyOk {
		return
	}

	if !strategy.typeSignatures.matches(args) {
		pgm.compileErrorString("invalid arguments for operation %s: %s", parts.op, args.String())
		return
	}

	instruction := strategy.build(pgm, op, args)
	pgm.instructions = append(pgm.instructions, instruction)

	return
}
