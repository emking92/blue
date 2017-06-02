package compiler

import (
	"regexp"
	"strconv"
)

type instructionBuildStrategy struct {
	typeSignatures instructionSignatureGroup
	build          instructionBuilder
}

type instructionBuilder func(pgm *programBuilder, op string, args []argument) Instruction

var opBuildStrategies = map[string]instructionBuildStrategy{
	"mov": instructionBuildStrategy{
		build: instructionBuilderMOV,
		typeSignatures: createSignatureGroup(
			[]string{"xx", "xx"},
			[]string{"xx", "[xx]"},
			[]string{"xx", "123"},
			[]string{"xx", "[123]"},
			[]string{"[123]", "xx"},
			[]string{"[xx]", "xx"},
			[]string{"[xx]", "123"},
		),
	},
	"var": instructionBuildStrategy{
		build: instructionBuilderVAR,
		typeSignatures: createSignatureGroup(
			[]string{"undefined"},
			[]string{"undefined", "xx"},
			[]string{"undefined", "123"},
		),
	},
	"jmp": instructionBuildStrategy{
		build: instructionBuilderJMP,
		typeSignatures: createSignatureGroup(
			[]string{"label"},
		),
	},
	"add": aluInstructionBuildStrategy,
	"sub": aluInstructionBuildStrategy,
	"mul": aluInstructionBuildStrategy,
	"div": aluInstructionBuildStrategy,
	"mod": aluInstructionBuildStrategy,
	"and": aluInstructionBuildStrategy,
	"or":  aluInstructionBuildStrategy,
	"xor": aluInstructionBuildStrategy,
	"sal": aluInstructionBuildStrategy,
	"sar": aluInstructionBuildStrategy,
}

var (
	aluInstructionBuildStrategy = instructionBuildStrategy{
		typeSignatures: createSignatureGroup(
			[]string{"xx", "xx", "xx"},
			[]string{"xx", "[123]", "xx"},
			[]string{"xx", "xx", "123"},
			[]string{"[123]", "xx", "xx"},
			[]string{"[123]", "xx", "123"},
		),
		build: instructionBuilderALU,
	}

	aluOperationValues = map[string]byte{
		"add": 1,
		"sub": 2,
		"mul": 3,
		"div": 4,
		"mod": 5,
		"and": 6,
		"or":  7,
		"xor": 8,
		"sal": 9,
		"sar": 10,
	}

	legalVariableNames   *regexp.Regexp
	illegalVariableNames *regexp.Regexp
)

func init() {
	legalVariableNames, _ = regexp.Compile(`^[a-zA-z_][a-zA-z0-9_]*$`)
	illegalVariableNames, _ = regexp.Compile(`(?i)^([a-z]x|[a-z0-9]p)$`)
}

func instructionBuilderMOV(pgm *programBuilder, op string, args []argument) Instruction {
	ins0 := args[0].build(pgm)
	ins1 := args[1].build(pgm)

	return joinInstructions(ins0, ins1)
}

func instructionBuilderALU(pgm *programBuilder, op string, args []argument) Instruction {
	aluByte, _ := aluOperationValues[op]
	insAlu := Instruction{Alu: aluByte}

	ins0 := args[0].build(pgm)
	ins1 := args[1].build(pgm)
	ins2 := args[2].build(pgm)

	return joinInstructions(insAlu, ins0, ins1, ins2)
}

func instructionBuilderVAR(pgm *programBuilder, op string, args []argument) Instruction {
	name := args[0].argStr
	if illegalVariableNames.MatchString(name) || !legalVariableNames.MatchString(name) {
		pgm.compileErrorString(`illegal variable name "%s"`, name)
		return Instruction{}
	}

	if pgm.variables.IsVarDefined(name) {
		pgm.compileErrorString(`variable already defined "%s"`, name)
		return Instruction{}
	}

	address := pgm.allocateWord()
	pgm.variables.CreateVariable(name, strconv.Itoa(address))

	toInstruction := Instruction{Wr: 1, Addr: address, Mar: 2, Mbr: 1}
	var fromInstruction Instruction

	if len(args) == 0 {
		fromInstruction = Instruction{Imm: 0, Cmux: 1}
	} else {
		fromInstruction = args[1].build(pgm)
	}

	return joinInstructions(toInstruction, fromInstruction)
}

func instructionBuilderJMP(pgm *programBuilder, op string, args []argument) Instruction {
	cond := Instruction{Cond: 3}
	addr := args[0].build(pgm)

	return joinInstructions(cond, addr)
}
