package compiler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type argumentType string

type argument struct {
	argStr    string
	index     int
	argType   argumentType
	buildFunc func(arg argument) (Instruction, error)
}

type argumentGroup []argument

const (
	argumentTypeUndefined        = argumentType("Undefined")
	argumentTypeImmediate        = argumentType("Immediate")
	argumentTypeImmediatePointer = argumentType("[Immediate]")
	argumentTypeRegister         = argumentType("Register")
	argumentTypeRegisterPointer  = argumentType("[Register]")
	argumentTypePort             = argumentType("Port")
	argumentTypeLabel            = argumentType("Label")
)

var (
	immediateArgumentRegex        *regexp.Regexp
	immediatePointerArgumentRegex *regexp.Regexp
	registerArgumentRegex         *regexp.Regexp
	registerPointerArgumentRegex  *regexp.Regexp
	registerMap                   map[string]byte
)

func init() {
	immediateArgumentRegex, _ = regexp.Compile(`^\d+$`)
	immediatePointerArgumentRegex, _ = regexp.Compile(`^\[(\d+)\]$`)
	registerArgumentRegex, _ = regexp.Compile(`^[a-z]x$`)
	registerPointerArgumentRegex, _ = regexp.Compile(`^\[([a-z]x)\]$`)
	registerMap = map[string]byte{
		"ax": 10,
		"bx": 11,
		"cx": 12,
		"dx": 13,
		"ex": 14,
		"fx": 15,
	}
}

func (arg argument) build(pgm *programBuilder) (ins Instruction) {
	ins, err := arg.buildFunc(arg)
	if err != nil {
		pgm.compileError(err)
	}
	return ins
}

func (group argumentGroup) types() (types []argumentType) {
	types = make([]argumentType, len(group))

	for i := range group {
		types[i] = group[i].argType
	}

	return types
}

func (group argumentGroup) String() string {
	strs := make([]string, len(group))
	for i, arg := range group {
		strs[i] = fmt.Sprintf("%s (type %s)", arg.argStr, arg.argType)
	}

	return fmt.Sprintf("{%s}", strings.Join(strs, ", "))
}

func (pgm programBuilder) parseInstructionArgument(argStr string, index int) (arg argument, err error) {
	arg = argument{
		argStr: argStr,
		index:  index,
	}

	if immediateArgumentRegex.MatchString(argStr) {
		arg.argType = argumentTypeImmediate
		arg.buildFunc = buildImmediateArgument

	} else if immediatePointerArgumentRegex.MatchString(argStr) {
		arg.argType = argumentTypeImmediatePointer
		arg.buildFunc = buildImmediatePointerArgument

	} else if registerArgumentRegex.MatchString(argStr) {
		arg.argType = argumentTypeRegister
		arg.buildFunc = buildRegisterArgument

	} else if registerPointerArgumentRegex.MatchString(argStr) {
		arg.argType = argumentTypeRegisterPointer
		arg.buildFunc = buildRegisterPointerArgument

	} else if pgm.isLabelDefined(argStr) {
		arg.argType = argumentTypeLabel
		arg.buildFunc = buildLabelArgument

	} else {
		arg.argType = argumentTypeUndefined
	}

	return arg, nil
}

func buildImmediateArgument(this argument) (ins Instruction, err error) {
	argInt, err := strconv.Atoi(this.argStr)
	if err == strconv.ErrRange {
		err = fmt.Errorf("immediate value out of int32 range: %s", this.argStr)
		return
	} else if err != nil {
		panic(err)
	}

	switch this.index {
	case 0:
		panic("Invalid argument")
	case 1:
		ins = Instruction{Imm: argInt, Cmux: 1}
	case 2:
		ins = Instruction{Imm: argInt, Bmux: 1}
	}

	return
}

func buildImmediatePointerArgument(this argument) (ins Instruction, err error) {
	return
}

func buildRegisterArgument(this argument) (ins Instruction, err error) {
	regByte, ok := registerMap[strings.ToLower(this.argStr)]
	if !ok {
		err = fmt.Errorf("undefined register: %s", this.argStr)
		return
	}

	switch this.index {
	case 0:
		ins = Instruction{C: regByte, Enc: 1}
	case 1:
		ins = Instruction{A: regByte, Amux: 0}
	case 2:
		ins = Instruction{B: regByte, Bmux: 0}
	}

	return
}

func buildRegisterPointerArgument(this argument) (ins Instruction, err error) {
	return
}

func buildLabelArgument(this argument) (ins Instruction, err error) {
	return
}

func buildPortArgument(this argument) (ins Instruction, err error) {
	return
}
