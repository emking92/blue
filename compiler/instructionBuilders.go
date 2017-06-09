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
			[]string{"[123]", "123"},
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
	"in": instructionBuildStrategy{
		build: instructionBuilderIO,
		typeSignatures: createSignatureGroup(
			[]string{"xx", "p_p"},
		),
	},
	"out": instructionBuildStrategy{
		build: instructionBuilderIO,
		typeSignatures: createSignatureGroup(
			[]string{"p_p", "xx"},
			[]string{"p_p", "123"},
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
	"je":  condJmpInstructionBuildStrategy,
	"jne": condJmpInstructionBuildStrategy,
	"jl":  condJmpInstructionBuildStrategy,
	"jle": condJmpInstructionBuildStrategy,
	"jg":  condJmpInstructionBuildStrategy,
	"jge": condJmpInstructionBuildStrategy,
}

var (
	aluInstructionBuildStrategy = instructionBuildStrategy{
		typeSignatures: createSignatureGroup(
			[]string{"xx", "xx", "xx"},
			[]string{"xx", "[123]", "xx"},
			[]string{"xx", "[123]", "123"},
			[]string{"xx", "xx", "123"},
			[]string{"[123]", "xx", "xx"},
			[]string{"[123]", "xx", "123"},
		),
		build: instructionBuilderALU,
	}
	condJmpInstructionBuildStrategy = instructionBuildStrategy{
		typeSignatures: createSignatureGroup(
			[]string{"label", "xx", "xx"},
			[]string{"label", "xx", "123"},
			[]string{"label", "[123]", "xx"},
			[]string{"label", "[123]", "123"},
		),
		build: instructionBuilderCondJMP,
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

	condOperationValues = map[string]byte{
		"je":  1,
		"jne": 2,
		"jl":  3,
		"jg":  4,
		"jle": 5,
		"jge": 6,
		"jmp": 7,
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
	jmpInstruction := Instruction{Cond: condOperationValues[op]}
	addr := args[0].build(pgm)

	return joinInstructions(jmpInstruction, addr)
}

func instructionBuilderIO(pgm *programBuilder, op string, args []argument) Instruction {
	ins0 := args[0].build(pgm)
	ins1 := args[1].build(pgm)

	return joinInstructions(ins0, ins1)
}

func instructionBuilderCondJMP(pgm *programBuilder, op string, args []argument) Instruction {
	jmpInstruction := Instruction{Alu: 2, Cond: condOperationValues[op]}
	addr := args[0].build(pgm)
	val1 := args[1].build(pgm)
	val2 := args[2].build(pgm)

	return joinInstructions(jmpInstruction, addr, val1, val2)
}
