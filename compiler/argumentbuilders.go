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
	buildFunc func(pgm *programBuilder, arg argument) (Instruction, error)
}

type argumentGroup []argument

const (
	argumentTypeUndefined          = argumentType("undefined")
	argumentTypeImmediate          = argumentType("immediate")
	argumentTypeImmediateReference = argumentType("immediate_reference")
	argumentTypeRegister           = argumentType("register")
	argumentTypeRegisterReference  = argumentType("register_reference")
	argumentTypePort               = argumentType("port")
	argumentTypeLabel              = argumentType("label")
)

var (
	addressOfRegex                  *regexp.Regexp
	immediateArgumentRegex          *regexp.Regexp
	immediateReferenceArgumentRegex *regexp.Regexp
	registerArgumentRegex           *regexp.Regexp
	registerReferenceArgumentRegex  *regexp.Regexp
	portArgumentRegex               *regexp.Regexp
	registerMap                     = map[string]byte{
		"eax": 10,
		"ebx": 11,
		"ecx": 12,
		"edx": 13,
		"eex": 14,
		"efx": 15,
	}
	portMap = map[string]int{
		"p_1": 1,
		"p_2": 2,
		"p_3": 3,
		"p_4": 4,
		"p_5": 5,
		"p_6": 6,
		"p_7": 7,
		"p_8": 8,
		"p_9": 9,
		"p_a": 10,
		"p_b": 11,
		"p_c": 12,
		"p_d": 13,
		"p_e": 14,
		"p_f": 15,
	}
)

func init() {
	addressOfRegex, _ = regexp.Compile(`&\[(.+)\]`)
	immediateArgumentRegex, _ = regexp.Compile(`^(-|0b|0x)?\d+$`)
	immediateReferenceArgumentRegex, _ = regexp.Compile(`^\[(-|0b|0x|0)?\d+\]$`)
	registerArgumentRegex, _ = regexp.Compile(`^e[a-z]x$`)
	registerReferenceArgumentRegex, _ = regexp.Compile(`^\[e[a-z]x\]$`)
	portArgumentRegex, _ = regexp.Compile(`^p_[a-z0-9]$`)
}

func (arg argument) build(pgm *programBuilder) (ins Instruction) {
	ins, err := arg.buildFunc(pgm, arg)
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
		strs[i] = fmt.Sprintf("%s (%s)", arg.argStr, arg.argType)
	}

	return fmt.Sprintf("{%s}", strings.Join(strs, ", "))
}

func (pgm programBuilder) parseInstructionArgument(argStr string, index int) (arg argument, err error) {
	argStr = pgm.variables.Expand(argStr)

	addressOf := addressOfRegex.FindStringSubmatch(argStr)
	if addressOf != nil {
		argStr = addressOf[1]
	}

	arg = argument{
		argStr: argStr,
		index:  index,
	}

	if immediateArgumentRegex.MatchString(argStr) {
		arg.argType = argumentTypeImmediate
		arg.buildFunc = buildImmediateArgument

	} else if immediateReferenceArgumentRegex.MatchString(argStr) {
		arg.argType = argumentTypeImmediateReference
		arg.buildFunc = buildImmediateReferenceArgument

	} else if registerArgumentRegex.MatchString(argStr) {
		arg.argType = argumentTypeRegister
		arg.buildFunc = buildRegisterArgument

	} else if registerReferenceArgumentRegex.MatchString(argStr) {
		arg.argType = argumentTypeRegisterReference
		arg.buildFunc = buildRegisterReferenceArgument

	} else if portArgumentRegex.MatchString(argStr) {
		arg.argType = argumentTypePort
		arg.buildFunc = buildPortArgument

	} else if pgm.isLabelDefined(argStr) {
		arg.argType = argumentTypeLabel
		arg.buildFunc = buildLabelArgument

	} else {
		arg.argType = argumentTypeUndefined
	}

	return arg, nil
}

func buildImmediateArgument(pgm *programBuilder, this argument) (ins Instruction, err error) {
	argInt, err := strconv.ParseInt(this.argStr, 0, 32)
	if err == strconv.ErrRange {
		err = fmt.Errorf("immediate value out of int32 range: %s", this.argStr)
		return
	} else if err != nil {
		panic(err)
	}

	argInt32 := int(argInt)

	switch this.index {
	case 0:
		panic("Invalid argument")
	case 1:
		ins = Instruction{Imm: argInt32, Cmux: 1}
	case 2:
		ins = Instruction{Imm: argInt32, Bmux: 1}
	}

	return
}

func buildImmediateReferenceArgument(pgm *programBuilder, this argument) (ins Instruction, err error) {
	argInt, err := strconv.ParseInt(strings.Trim(this.argStr, "[]"), 0, 32)
	if err == strconv.ErrRange {
		err = fmt.Errorf("reference address out of int32 range: %s", this.argStr)
		return
	} else if err != nil {
		panic(err)
	}

	argInt32 := int(argInt)

	switch this.index {
	case 0:
		ins = Instruction{Addr: argInt32, Wr: 1, Mar: 2, Mbr: 1}
	case 1:
		ins = Instruction{Addr: argInt32, Rd: 1, Mar: 2, Amux: 1}
	case 2:
		panic("Invalid argument")
	}

	return
}

func buildRegisterArgument(pgm *programBuilder, this argument) (ins Instruction, err error) {
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

func buildRegisterReferenceArgument(pgm *programBuilder, this argument) (ins Instruction, err error) {
	reg := strings.ToLower(strings.Trim(this.argStr, "[]"))
	regByte, ok := registerMap[reg]
	if !ok {
		err = fmt.Errorf("undefined register reference: %s", this.argStr)
		return
	}

	switch this.index {
	case 0:
		ins = Instruction{B: regByte, Wr: 1, Mar: 1, Mbr: 1}
	case 1:
		ins = Instruction{B: regByte, Rd: 1, Mar: 1, Amux: 1}
	case 2:
		panic("Invalid argument")
	}

	return
}

func buildLabelArgument(pgm *programBuilder, this argument) (ins Instruction, err error) {
	if this.index != 0 {
		panic("Invalid argument")
	}

	label := this.argStr
	labelLine, ok := pgm.getLabelLine(label)
	if !ok {
		err = fmt.Errorf(`undefined label "%s"`, label)
		ins = Instruction{}
		return
	}

	ins = Instruction{Bran: labelLine}
	return
}

func buildPortArgument(pgm *programBuilder, this argument) (ins Instruction, err error) {
	portAddress, ok := portMap[strings.ToLower(this.argStr)]
	if !ok {
		err = fmt.Errorf("undefined port: %s", this.argStr)
		return
	}

	switch this.index {
	case 0:
		ins = Instruction{Addr: portAddress, Wr: 1, Mar: 2, Mbr: 1}
	case 1:
		ins = Instruction{Addr: portAddress, Rd: 1, Mar: 2, Amux: 1}
	case 2:
		panic("Invalid argment")
	}

	return
}
